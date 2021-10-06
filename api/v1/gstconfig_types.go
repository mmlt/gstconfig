/*
Copyright 2021.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GSTConfigSpec defines the desired state of GSTConfig
type GSTConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of GSTConfig. Edit gstconfig_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// GSTConfigStatus defines the observed state of GSTConfig
type GSTConfigStatus struct {
	// Conditions are a synopsis of the StepStates.
	// +optional
	Conditions []GSTConfigCondition `json:"conditions,omitempty"`
}

// GSTConfigCondition provides a synopsis of the current state.
// See KEP sig-api-machinery/1623-standardize-conditions is going to introduce it as k8s.io/apimachinery/pkg/apis/meta/v1
type GSTConfigCondition struct {
	// Type of condition in CamelCase.
	// +required
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status of the condition, one of True, False, Unknown.
	// +required
	Status metav1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// Last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition in CamelCase.
	// +required
	Reason GSTConfigConditionReason `json:"reason" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// GSTConfigConditionReason is the reason for the condition change.
type GSTConfigConditionReason string

const (
	ReasonRunning GSTConfigConditionReason = "Running"
	ReasonReady   GSTConfigConditionReason = "Ready"
	ReasonFailed  GSTConfigConditionReason = "Failed"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GSTConfig is the Schema for the gstconfigs API
type GSTConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GSTConfigSpec   `json:"spec,omitempty"`
	Status GSTConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GSTConfigList contains a list of GSTConfig
type GSTConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GSTConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GSTConfig{}, &GSTConfigList{})
}
