/*
Copyright 2021 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

// +k8s:deepcopy-gen=package
// +groupName=rke.cattle.io
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RKEBootstrapList is a list of RKEBootstrap resources
type RKEBootstrapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RKEBootstrap `json:"items"`
}

func NewRKEBootstrap(namespace, name string, obj RKEBootstrap) *RKEBootstrap {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("RKEBootstrap").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RKEBootstrapTemplateList is a list of RKEBootstrapTemplate resources
type RKEBootstrapTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RKEBootstrapTemplate `json:"items"`
}

func NewRKEBootstrapTemplate(namespace, name string, obj RKEBootstrapTemplate) *RKEBootstrapTemplate {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("RKEBootstrapTemplate").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RKEClusterList is a list of RKECluster resources
type RKEClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RKECluster `json:"items"`
}

func NewRKECluster(namespace, name string, obj RKECluster) *RKECluster {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("RKECluster").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RKEControlPlanList is a list of RKEControlPlan resources
type RKEControlPlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []RKEControlPlan `json:"items"`
}

func NewRKEControlPlan(namespace, name string, obj RKEControlPlan) *RKEControlPlan {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("RKEControlPlan").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UnmanagedMachineList is a list of UnmanagedMachine resources
type UnmanagedMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []UnmanagedMachine `json:"items"`
}

func NewUnmanagedMachine(namespace, name string, obj UnmanagedMachine) *UnmanagedMachine {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("UnmanagedMachine").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}
