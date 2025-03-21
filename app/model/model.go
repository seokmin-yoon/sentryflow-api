package model

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APILog struct {
	ID        uint64 `json:"id"`
	TimeStamp string `json:"timeStamp"`

	SrcNamespace string `json:"srcNamespace"`
	SrcName      string `json:"srcName"`
	//	SrcLabel     map[string]string `json:"srcLabel"`
	SrcType string `json:"srcType"`
	SrcIP   string `json:"srcIP"`
	SrcPort string `json:"srcPort"`

	DstNamespace string `json:"dstNamespace"`
	DstName      string `json:"dstName"`
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
	Pods       []Pod       `json:"pods"`
}

type Pod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PodIP       string `json:"pod_ip"`
	ClusterName string `json:"cluster_name"`
	Namespace   string `json:"namespace"`
	QosClass    string `json:"qos_class"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
}

type Namespace struct {
	Name  string `json:"name"`
	Phase string `json:"phase"` // Active or Terminating
}
