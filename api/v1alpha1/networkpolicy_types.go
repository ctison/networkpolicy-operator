package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NetworkPolicySpec defines the desired state of NetworkPolicy
type NetworkPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NetworkPolicy. Edit NetworkPolicy_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// NetworkPolicyStatus defines the observed state of NetworkPolicy
type NetworkPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// NetworkPolicy is the Schema for the networkpolicies API
type NetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkPolicySpec   `json:"spec,omitempty"`
	Status NetworkPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkPolicyList contains a list of NetworkPolicy
type NetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkPolicy{}, &NetworkPolicyList{})
}
