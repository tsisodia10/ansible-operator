// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybook) DeepCopyInto(out *AnsiblePlaybook) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybook.
func (in *AnsiblePlaybook) DeepCopy() *AnsiblePlaybook {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybook)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnsiblePlaybook) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookList) DeepCopyInto(out *AnsiblePlaybookList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AnsiblePlaybook, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookList.
func (in *AnsiblePlaybookList) DeepCopy() *AnsiblePlaybookList {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnsiblePlaybookList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookRun) DeepCopyInto(out *AnsiblePlaybookRun) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookRun.
func (in *AnsiblePlaybookRun) DeepCopy() *AnsiblePlaybookRun {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnsiblePlaybookRun) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookRunList) DeepCopyInto(out *AnsiblePlaybookRunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AnsiblePlaybookRun, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookRunList.
func (in *AnsiblePlaybookRunList) DeepCopy() *AnsiblePlaybookRunList {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookRunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AnsiblePlaybookRunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookRunSpec) DeepCopyInto(out *AnsiblePlaybookRunSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookRunSpec.
func (in *AnsiblePlaybookRunSpec) DeepCopy() *AnsiblePlaybookRunSpec {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookRunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookRunStatus) DeepCopyInto(out *AnsiblePlaybookRunStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookRunStatus.
func (in *AnsiblePlaybookRunStatus) DeepCopy() *AnsiblePlaybookRunStatus {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookSpec) DeepCopyInto(out *AnsiblePlaybookSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookSpec.
func (in *AnsiblePlaybookSpec) DeepCopy() *AnsiblePlaybookSpec {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AnsiblePlaybookStatus) DeepCopyInto(out *AnsiblePlaybookStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AnsiblePlaybookStatus.
func (in *AnsiblePlaybookStatus) DeepCopy() *AnsiblePlaybookStatus {
	if in == nil {
		return nil
	}
	out := new(AnsiblePlaybookStatus)
	in.DeepCopyInto(out)
	return out
}
