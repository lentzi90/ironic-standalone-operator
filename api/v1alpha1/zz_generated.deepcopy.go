//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DHCP) DeepCopyInto(out *DHCP) {
	*out = *in
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ignore != nil {
		in, out := &in.Ignore, &out.Ignore
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DHCP.
func (in *DHCP) DeepCopy() *DHCP {
	if in == nil {
		return nil
	}
	out := new(DHCP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Images) DeepCopyInto(out *Images) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Images.
func (in *Images) DeepCopy() *Images {
	if in == nil {
		return nil
	}
	out := new(Images)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Inspection) DeepCopyInto(out *Inspection) {
	*out = *in
	if in.Collectors != nil {
		in, out := &in.Collectors, &out.Collectors
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VLANInterfaces != nil {
		in, out := &in.VLANInterfaces, &out.VLANInterfaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Inspection.
func (in *Inspection) DeepCopy() *Inspection {
	if in == nil {
		return nil
	}
	out := new(Inspection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ironic) DeepCopyInto(out *Ironic) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ironic.
func (in *Ironic) DeepCopy() *Ironic {
	if in == nil {
		return nil
	}
	out := new(Ironic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ironic) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicDatabase) DeepCopyInto(out *IronicDatabase) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicDatabase.
func (in *IronicDatabase) DeepCopy() *IronicDatabase {
	if in == nil {
		return nil
	}
	out := new(IronicDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IronicDatabase) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicDatabaseList) DeepCopyInto(out *IronicDatabaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IronicDatabase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicDatabaseList.
func (in *IronicDatabaseList) DeepCopy() *IronicDatabaseList {
	if in == nil {
		return nil
	}
	out := new(IronicDatabaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IronicDatabaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicDatabaseSpec) DeepCopyInto(out *IronicDatabaseSpec) {
	*out = *in
	out.CredentialsRef = in.CredentialsRef
	out.TLSRef = in.TLSRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicDatabaseSpec.
func (in *IronicDatabaseSpec) DeepCopy() *IronicDatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(IronicDatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicDatabaseStatus) DeepCopyInto(out *IronicDatabaseStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicDatabaseStatus.
func (in *IronicDatabaseStatus) DeepCopy() *IronicDatabaseStatus {
	if in == nil {
		return nil
	}
	out := new(IronicDatabaseStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicList) DeepCopyInto(out *IronicList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ironic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicList.
func (in *IronicList) DeepCopy() *IronicList {
	if in == nil {
		return nil
	}
	out := new(IronicList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IronicList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicSpec) DeepCopyInto(out *IronicSpec) {
	*out = *in
	out.CredentialsRef = in.CredentialsRef
	out.DatabaseRef = in.DatabaseRef
	out.Images = in.Images
	in.Inspection.DeepCopyInto(&out.Inspection)
	in.Networking.DeepCopyInto(&out.Networking)
	out.TLSRef = in.TLSRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicSpec.
func (in *IronicSpec) DeepCopy() *IronicSpec {
	if in == nil {
		return nil
	}
	out := new(IronicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IronicStatus) DeepCopyInto(out *IronicStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IronicStatus.
func (in *IronicStatus) DeepCopy() *IronicStatus {
	if in == nil {
		return nil
	}
	out := new(IronicStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	if in.DHCP != nil {
		in, out := &in.DHCP, &out.DHCP
		*out = new(DHCP)
		(*in).DeepCopyInto(*out)
	}
	if in.MACAddresses != nil {
		in, out := &in.MACAddresses, &out.MACAddresses
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Networking.
func (in *Networking) DeepCopy() *Networking {
	if in == nil {
		return nil
	}
	out := new(Networking)
	in.DeepCopyInto(out)
	return out
}
