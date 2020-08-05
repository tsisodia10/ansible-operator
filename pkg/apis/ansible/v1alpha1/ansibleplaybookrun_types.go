package v1alpha1

import (
	"context"

	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	Pending   = "Pending"
	Preparing = "Preparing"
	Active    = "Active"
	Cleaning  = "Cleaning"
	Finished  = "Finished"
	Failed    = "Failed"
)

//https://mojo.redhat.com/docs/DOC-1207144
// AnsiblePlaybookRunSpec defines the desired state of AnsiblePlaybookRun
type AnsiblePlaybookRunSpec struct {
	AnsiblePlaybookRef *kapi.ObjectReference `json:"ansiblePlaybook,omitempty"`
	Inventory          string                `json:"inventory,omitempty"`
	Password           string                `json:"password,omitempty"`
	SSHPrivateKey      string                `json:"sshPrivateKey,omitempty"`
	ExtraVars          string                `json:"extraVars,omitempty"`
}

// AnsiblePlaybookRunStatus defines the observed state of AnsiblePlaybookRun
type AnsiblePlaybookRunStatus struct {
	Status string `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybookRun is the Schema for the ansibleplaybookruns API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ansibleplaybookruns,scope=Namespaced
type AnsiblePlaybookRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnsiblePlaybookRunSpec   `json:"spec,omitempty"`
	Status AnsiblePlaybookRunStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybookRunList contains a list of AnsiblePlaybookRun
type AnsiblePlaybookRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnsiblePlaybookRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnsiblePlaybookRun{}, &AnsiblePlaybookRunList{})
}

func (apr *AnsiblePlaybookRun) GetAnsiblePlaybook(client k8sclient.Client) (*AnsiblePlaybook, error) {

	//fmt.Printf("%+v\n", apr)
	if apr == nil {
		return nil, nil
	}
	object := AnsiblePlaybook{}
	// fmt.Printf("%+s\n", apr.Spec.AnsiblePlaybookRef.Namespace)
	// fmt.Printf("%+s\n", apr.Spec.AnsiblePlaybookRef.Name)
	err := client.Get(
		context.TODO(),
		types.NamespacedName{
			Namespace: apr.Spec.AnsiblePlaybookRef.Namespace,
			Name:      apr.Spec.AnsiblePlaybookRef.Name,
		},
		&object)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &object, err
}
