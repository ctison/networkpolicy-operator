package v1alpha1

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DefaultResolveEverySeconds uint64 = 300

// NetworkPolicySpec defines the desired state of NetworkPolicy
type NetworkPolicySpec struct {
	// Include native NetworkPolicy spec
	PodSelector metav1.LabelSelector                    `json:"podSelector"`
	Ingress     []networkingv1.NetworkPolicyIngressRule `json:"ingress,omitempty"`
	Egress      []NetworkPolicyEgressRule               `json:"egress,omitempty"`
	PolicyTypes []networkingv1.PolicyType               `json:"policyTypes,omitempty"`

	// Time interval in seconds to periodically resolve hostnames to IPs from DNS.
	// Defaults to 300 (5m).
	ResolveEverySeconds *uint64 `json:"resolveEverySeconds,omitempty"`
}

// NetworkPolicyStatus defines the observed state of NetworkPolicy
type NetworkPolicyStatus struct {
	NetworkPolicyName *string `json:"networkPolicyName,omitempty"`
}

type NetworkPolicyEgressRule struct {
	Ports []networkingv1.NetworkPolicyPort `json:"ports,omitempty"`
	To    []NetworkPolicyPeer              `json:"to,omitempty"`
}

type NetworkPolicyPeer struct {
	networkingv1.NetworkPolicyPeer `json:",inline"`
	Domain                         *string `json:"domain,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkPolicy is the Schema for the networkpolicies API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn name=Interval type= JSONPath=.spec.
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
