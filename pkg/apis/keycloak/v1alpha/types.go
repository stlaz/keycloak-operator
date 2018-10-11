package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeycloakOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []KeycloakOperator `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type KeycloakOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              KeycloakOperatorSpec   `json:"spec"`
	Status            KeycloakOperatorStatus `json:"status,omitempty"`
}

type KeycloakOperatorSpec struct {
	// Fill me
}
type KeycloakOperatorStatus struct {
	// Fill me
}
