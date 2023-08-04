package ironic

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

const (
	ironicAppName = "ironic-service"

	ironicPortName    = "ironic-api"
	imagesPortName    = "image-svc"
	imagesTLSPortName = "image-svc-tls"

	ironicUser     = 997
	ironicGroup    = 994
	inspectorUser  = 996
	inspectorGroup = 993

	authDir   = "/auth"
	certsDir  = "/certs"
	sharedDir = "/shared"
)

func buildCommonEnvVars(ironic *metal3api.Ironic) []corev1.EnvVar {
	result := []corev1.EnvVar{
		{
			Name:  "RESTART_CONTAINER_CERTIFICATE_UPDATED",
			Value: "true",
		},
		{
			Name:  "IRONIC_LISTEN_PORT",
			Value: strconv.Itoa(int(ironic.Spec.APIPort)),
		},
		{
			Name:  "HTTP_PORT",
			Value: strconv.Itoa(int(ironic.Spec.ImageServerPort)),
		},
		{
			Name:  "LISTEN_ALL_INTERFACES",
			Value: strconv.FormatBool(!ironic.Spec.Networking.BindInterface),
		},
	}

	networkingProvided := false
	if ironic.Spec.Networking.IPAddress != "" {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_IP",
				Value: ironic.Spec.Networking.IPAddress,
			},
		)
		networkingProvided = true
	}
	if ironic.Spec.Networking.Interface != "" {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_INTERFACE",
				Value: ironic.Spec.Networking.Interface,
			},
		)
		networkingProvided = true
	}
	if len(ironic.Spec.Networking.MACAddresses) > 0 {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_MACS",
				Value: strings.Join(ironic.Spec.Networking.MACAddresses, ","),
			},
		)
		networkingProvided = true
	}
	if !networkingProvided {
		result = append(result,
			corev1.EnvVar{
				Name: "PROVISIONING_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.hostIP",
					},
				},
			},
		)
	}

	if ironic.Spec.TLSSecretName != "" {
		// Ironic and Inspector will listen on a Unix socket, httpd will be responsible for serving HTTPS.
		result = append(result, []corev1.EnvVar{
			{
				Name:  "IRONIC_PRIVATE_PORT",
				Value: "unix",
			},
			{
				Name:  "IRONIC_INSPECTOR_PRIVATE_PORT",
				Value: "unix",
			},
			{
				Name:  "IRONIC_REVERSE_PROXY_SETUP",
				Value: "true",
			},
			{
				Name:  "INSPECTOR_REVERSE_PROXY_SETUP",
				Value: "true",
			},
		}...)

		if !ironic.Spec.DisableVirtualMediaTLS {
			result = append(result,
				corev1.EnvVar{
					Name:  "VMEDIA_TLS_PORT",
					Value: strconv.Itoa(int(ironic.Spec.ImageServerTLSPort)),
				},
			)
		}
	}

	if ironic.Spec.RamdiskSSHKey != "" {
		result = append(result,
			corev1.EnvVar{
				Name:  "IRONIC_RAMDISK_SSH_KEY",
				Value: strings.Trim(ironic.Spec.RamdiskSSHKey, " \t\n\r"),
			},
		)
	}

	return result
}

func buildHtpasswdVars(htpasswd string) (result []corev1.EnvVar) {
	if htpasswd != "" {
		result = append(result, []corev1.EnvVar{
			{
				Name:  "IRONIC_HTPASSWD",
				Value: htpasswd,
			},
			{
				Name:  "INSPECTOR_HTPASSWD",
				Value: htpasswd,
			},
		}...)
	}
	return
}

func buildIronicEnvVars(ironic *metal3api.Ironic, db *metal3api.IronicDatabase, htpasswd string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)
	result = append(result, []corev1.EnvVar{
		{
			Name:  "IRONIC_USE_MARIADB",
			Value: strconv.FormatBool(db != nil),
		},
		{
			Name:  "IRONIC_EXPOSE_JSON_RPC",
			Value: strconv.FormatBool(ironic.Spec.Distributed),
		},
		// TODO(dtantsur): try to get rid of these eventually (especially once inspector is gone)
		{
			Name:  "IRONIC_INSECURE",
			Value: "true",
		},
		{
			Name:  "IRONIC_INSPECTOR_INSECURE",
			Value: "true",
		},
	}...)

	if db != nil {
		result = append(result, databasePasswordEnvVar(db))
	}

	if ironic.Spec.TLSSecretName == "" {
		result = append(result, buildHtpasswdVars(htpasswd)...)
	}

	return result
}

func buildHttpdEnvVars(ironic *metal3api.Ironic, db *metal3api.IronicDatabase, htpasswd string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)

	if ironic.Spec.TLSSecretName != "" {
		result = append(result, buildHtpasswdVars(htpasswd)...)
	}

	return result
}

func buildIronicVolumesAndMounts(ironic *metal3api.Ironic, db *metal3api.IronicDatabase) (volumes []corev1.Volume, mounts []corev1.VolumeMount) {
	volumes = []corev1.Volume{
		{
			Name: "ironic-shared",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "ironic-auth",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ironic.Spec.APISecretName,
				},
			},
		},
		{
			Name: "inspector-auth",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ironic.Spec.APISecretName,
				},
			},
		},
	}
	mounts = []corev1.VolumeMount{
		{
			Name:      "ironic-shared",
			MountPath: sharedDir,
		},
		{
			Name:      "ironic-auth",
			MountPath: authDir + "/ironic",
		},
		{
			Name:      "inspector-auth",
			MountPath: authDir + "/ironic-inspector",
		},
	}

	if ironic.Spec.TLSSecretName != "" {
		volumes = append(volumes,
			corev1.Volume{
				Name: "cert-ironic",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: ironic.Spec.TLSSecretName,
					},
				},
			},
		)
		mounts = append(mounts,
			corev1.VolumeMount{
				Name:      "cert-ironic",
				MountPath: certsDir + "/ironic",
				ReadOnly:  true,
			},
			// NOTE(dtantsur): in theory, inspector can use different TLS settings, but we're getting rid of inspector,
			// so I don't want any API bits to rely on its presence.
			corev1.VolumeMount{
				Name:      "cert-ironic",
				MountPath: certsDir + "/ironic-inspector",
				ReadOnly:  true,
			},
		)
		if !ironic.Spec.DisableVirtualMediaTLS {
			mounts = append(mounts,
				corev1.VolumeMount{
					Name:      "cert-vmedia",
					MountPath: certsDir + "/vmedia",
					ReadOnly:  true,
				},
			)
		}
	}

	if db != nil && db.Spec.TLSSecretName != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "cert-mariadb",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: db.Spec.TLSSecretName,
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "cert-mariadb",
			MountPath: certsDir + "/mariadb",
		})
	}

	return
}

func buildIronicHttpdPorts(ironic *metal3api.Ironic) (ironicPorts []corev1.ContainerPort, httpdPorts []corev1.ContainerPort) {
	httpdPorts = []corev1.ContainerPort{
		{
			Name:          imagesPortName,
			ContainerPort: ironic.Spec.ImageServerPort,
		},
	}

	apiPort := corev1.ContainerPort{
		Name:          ironicPortName,
		ContainerPort: ironic.Spec.APIPort,
	}

	if ironic.Spec.TLSSecretName == "" {
		ironicPorts = append(ironicPorts, apiPort)
	} else {
		httpdPorts = append(httpdPorts, apiPort)
		if !ironic.Spec.DisableVirtualMediaTLS {
			httpdPorts = append(httpdPorts, corev1.ContainerPort{
				Name:          imagesTLSPortName,
				ContainerPort: ironic.Spec.ImageServerTLSPort,
			})
		}
	}

	return
}

func newIronicPodTemplate(ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (corev1.PodTemplateSpec, error) {
	var htpasswd string
	if apiSecret != nil {
		var err error
		htpasswd, err = htpasswdFromSecret(apiSecret)
		if err != nil {
			return corev1.PodTemplateSpec{}, fmt.Errorf("cannot generate htpasswd for the API credentials secret: %w", err)
		}
	}

	volumes, mounts := buildIronicVolumesAndMounts(ironic, db)
	sharedVolumeMount := mounts[0]
	initContainers := []corev1.Container{
		{
			Name:         "ipa-downloader",
			Image:        ironic.Spec.RamdiskDownloaderImage,
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
			},
			// FIXME(dtantsur): this should be the default in the image:
			// https://github.com/metal3-io/ironic-ipa-downloader/pull/43
			Command: []string{"/usr/local/bin/get-resource.sh"},
		},
	}

	ironicPorts, httpdPorts := buildIronicHttpdPorts(ironic)

	containers := []corev1.Container{
		{
			Name:    "ironic",
			Image:   ironic.Spec.Image,
			Command: []string{"/bin/runironic"},
			// TODO(dtantsur): livenessProbe+readinessProbe
			Env:          buildIronicEnvVars(ironic, db, htpasswd),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
			},
			Ports: ironicPorts,
		},
		{
			Name:    "httpd",
			Image:   ironic.Spec.Image,
			Command: []string{"/bin/runhttpd"},
			// TODO(dtantsur): livenessProbe+readinessProbe
			Env:          buildHttpdEnvVars(ironic, db, htpasswd),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
			},
			Ports: httpdPorts,
		},
		{
			Name:         "ramdisk-logs",
			Image:        ironic.Spec.Image,
			Command:      []string{"/bin/runlogwatch.sh"},
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
			},
		},
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{metal3api.IronicOperatorLabel: ironicAppName},
		},
		Spec: corev1.PodSpec{
			Containers:     containers,
			InitContainers: initContainers,
			Volumes:        volumes,
			// Ironic needs to be accessed by external machines
			HostNetwork: true,
		},
	}, nil
}

func ensureIronicDaemonSet(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, err error) {
	if db == nil {
		return metal3api.IronicStatusProgressing, errors.New("database is required for a distributed deployment")
	}
	template, err := newIronicPodTemplate(ironic, db, apiSecret)
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}

	deploy := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	_, err = controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, deploy, func() error {
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			cctx.Logger.Info("creating a new ironic daemon set")
			matchLabels := map[string]string{metal3api.IronicOperatorLabel: ironicAppName}
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: matchLabels,
			}
		}
		deploy.Spec.Template = template
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}
	return getDaemonSetStatus(deploy)
}

func ensureIronicDeployment(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, err error) {
	template, err := newIronicPodTemplate(ironic, db, apiSecret)
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	_, err = controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, deploy, func() error {
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			cctx.Logger.Info("creating a new ironic deployment")
			matchLabels := map[string]string{metal3api.IronicOperatorLabel: ironicAppName}
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: matchLabels,
			}
			deploy.Spec.Replicas = pointer.Int32(1)
		}
		deploy.Spec.Template = template
		// We cannot run two copies of Ironic in parallel
		deploy.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		}
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}
	return getDeploymentStatus(deploy)
}

func ensureIronicService(cctx ControllerContext, ironic *metal3api.Ironic) (metal3api.IronicStatusConditionType, []string, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	exposedPort := 80
	proto := "http"
	if ironic.Spec.TLSSecretName != "" {
		exposedPort = 443
		proto = "https"
	}
	_, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, service, func() error {
		if service.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new ironic service")
			service.ObjectMeta.Labels = make(map[string]string)
		}
		service.ObjectMeta.Labels[metal3api.IronicOperatorLabel] = databaseAppName

		service.Spec.Selector = map[string]string{metal3api.IronicOperatorLabel: ironicAppName}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(exposedPort),
				TargetPort: intstr.FromString(ironicPortName),
			},
		}
		service.Spec.Type = corev1.ServiceTypeClusterIP

		return controllerutil.SetControllerReference(ironic, service, cctx.Scheme)
	})
	if err != nil || len(service.Spec.ClusterIPs) == 0 {
		return metal3api.IronicStatusProgressing, nil, err
	}

	return metal3api.IronicStatusAvailable, buildEndpoints(service.Spec.ClusterIPs, exposedPort, proto), nil
}

func removeIronicDaemonSet(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil // TODO(dtantsur): migration, cannot use ownership
}

func removeIronicDeployment(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil // TODO(dtantsur): migration, cannot use ownership
}

// EnsureIronic deploys Ironic either as a Deployment or as a DaemonSet.
func EnsureIronic(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, endpoints []string, err error) {
	if ironic.Spec.Distributed {
		err = removeIronicDeployment(cctx, ironic)
		if err != nil {
			return
		}
		status, err = ensureIronicDaemonSet(cctx, ironic, db, apiSecret)
	} else {
		err = removeIronicDaemonSet(cctx, ironic)
		if err != nil {
			return
		}
		status, err = ensureIronicDeployment(cctx, ironic, db, apiSecret)
	}

	if err != nil {
		return
	}

	// Let the service be created while Ironic is being deployed, but do
	// not report overall success until both are done.
	serviceStatus, endpoints, err := ensureIronicService(cctx, ironic)
	if err != nil || serviceStatus != metal3api.IronicStatusAvailable {
		return serviceStatus, nil, err
	}

	return
}

// RemoveIronic removes all bits of the Ironic deployment.
func RemoveIronic(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil // rely on ownership-based clean up
}
