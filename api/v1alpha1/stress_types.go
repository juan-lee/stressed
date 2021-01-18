/*
Copyright 2021 Juan-Lee Pang.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StressSpec defines the desired state of Stress
type StressSpec struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Image defines the container image the test will use.
	// +optional
	Image string `json:"image,omitempty"`

	// JobFile contains a stress-ng jobfile.
	// More info: `man stress-ng` for possible parameters.
	// +optional
	JobFile string `json:"jobFile,omitempty"`
}

// StressStatus defines the observed state of Stress
type StressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Stress is the Schema for the stresses API
type Stress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StressSpec   `json:"spec,omitempty"`
	Status StressStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StressList contains a list of Stress
type StressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stress{}, &StressList{})
}
