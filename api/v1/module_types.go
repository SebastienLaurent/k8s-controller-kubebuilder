/*
Copyright 2022.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ModuleSpec defines the desired state of Module

type SubModuleSpec struct {
	Val1 string `json:"val1"`
	Val2 string `json:"val2"`
}

type ModuleSpec struct {
	// +kubebuilder:validation:MaxLength=15
	Cu string `json:"cu"`

	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:MinLength=1
	Module string `json:"module"`

	Sidecar corev1.PodSpec `json:"sidecar"`

	//+listType=map
	//+listMapKey=val1
	Lst []SubModuleSpec `json:"lst"`
}

// ModuleStatus defines the observed state of Module
type ModuleStatus struct {

	// A list of pointers to currently running jobs.
	// +optional
	Sidecar *corev1.ObjectReference `json:"sidecar,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=mod
//+kubebuilder:printcolumn:name="CU",type=string,JSONPath=`.spec.cu`

// Module is the Schema for the modules API
type Module struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModuleSpec   `json:"spec,omitempty"`
	Status ModuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ModuleList contains a list of Module
type ModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Module `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Module{}, &ModuleList{})
}
