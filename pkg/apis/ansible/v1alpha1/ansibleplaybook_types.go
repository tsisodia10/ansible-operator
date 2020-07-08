package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AnsiblePlaybookSpec defines the desired state of AnsiblePlaybook
type AnsiblePlaybookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	RepositoryType string `json:"repository,omitempty"`
	RepositoryURL  string `json:"url,omitempty"`
	PlaybookName   string `json:"playbookName,omitempty"`
}

// AnsiblePlaybookStatus defines the observed state of AnsiblePlaybook
type AnsiblePlaybookStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybook is the Schema for the ansibleplaybooks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=ansibleplaybooks,scope=Namespaced
type AnsiblePlaybook struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AnsiblePlaybookSpec   `json:"spec,omitempty"`
	Status AnsiblePlaybookStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AnsiblePlaybookList contains a list of AnsiblePlaybook
type AnsiblePlaybookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AnsiblePlaybook `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AnsiblePlaybook{}, &AnsiblePlaybookList{})
}
