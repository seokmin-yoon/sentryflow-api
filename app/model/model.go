package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APILog struct {
	ID        uint64 `json:"id"`
	TimeStamp string `json:"timeStamp"`

	SrcNamespace string            `json:"srcNamespace"`
	SrcName      string            `json:"srcName"`
//	SrcLabel     map[string]string `json:"srcLabel"`
	SrcType string `json:"srcType"`
	SrcIP   string `json:"srcIP"`
	SrcPort string `json:"srcPort"`

	DstNamespace string            `json:"dstNamespace"`
	DstName      string            `json:"dstName"`
//	DstLabel     map[string]string `json:"dstLabel"`
	DstType string `json:"dstType"`
	DstIP   string `json:"dstIP"`
	DstPort string `json:"dstPort"`

//	Protocol     string `json:"protocol"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	ResponseCode int32  `json:"responseCode"`
}

type APIMetrics struct {
	PerAPICounts map[string]uint64 `json:"perAPICounts"`
}

type MetricValue struct {
	Value map[string]string `json:"value"`
}

type EnvoyMetrics struct {
	TimeStamp string                 `json:"timeStamp"`
	Namespace string                 `json:"namespace"`
	Name      string                 `json:"name"`
	IPAddress string                 `json:"ipAddress"`
	Labels    map[string]string      `json:"labels"`
	Metrics   map[string]MetricValue `json:"metrics"`
}

type Cluster struct {
	Name       string      `json:"name"`
	Namespaces []Namespace `json:"namespaces"`
	Pods	   []Pod	   `json:"pods"`
}

type Pod struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   v1.PodSpec   `json:"spec,omitempty"`
    Status v1.PodStatus `json:"status,omitempty"`
}

type Namespace struct {
	Name  string `json:"name"`
	Phase string `json:"phase"` // Active or Terminating
}
