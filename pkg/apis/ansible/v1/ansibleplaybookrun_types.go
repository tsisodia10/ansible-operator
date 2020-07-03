package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AnsiblePlaybookRunSpec defines the desired state of AnsiblePlaybookRun
type AnsiblePlaybookRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// AnsiblePlaybookRunStatus defines the observed state of AnsiblePlaybookRun
type AnsiblePlaybookRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
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
