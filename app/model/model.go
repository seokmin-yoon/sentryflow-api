package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APILog struct {
	ID        int64  `json:"id"`
	TimeStamp string `json:"timeStamp"`

	SrcCluster   string `json:"srcCluster"`
	SrcNamespace string `json:"srcNamespace"`
	SrcName      string `json:"srcName"`
	//	SrcLabel     map[string]string `json:"srcLabel"`
	SrcType string `json:"srcType"`
	SrcIP   string `json:"srcIP"`
	SrcPort string `json:"srcPort"`

	DstCluster   string `json:"dstCluster"`
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
	Name       string   `json:"name"`
	Namespaces []string `json:"namespaces"`
}

type Pod struct {
	ObjectID          primitive.ObjectID `bson:"_id" json:"id"`
	Cluster           string             `json:"cluster"`
	Namespace         string             `json:"namespace"`
	Name              string             `json:"name"`
	Nodename          string             `json:"nodename"`
	PodIP             string             `json:"podip"`
	Status            string             `json:"status"`
	CreationTimestamp string             `json:"creationTimestamp"`
}

type Service struct {
	ObjectID        primitive.ObjectID `bson:"_id" json:"id"`
	Cluster         string             `json:"cluster"`
	Namespace       string             `json:"namespace"`
	Name            string             `json:"name"`
	Type            string             `json:"type"`
	ClusterIP       string             `json:"clusterip"`
	Ports           []Port             `json:"ports"`
	LoadBalancerIPs []string           `json:"loadbalancerips"`
}

type Port struct {
	Port       int    `json:"port"`
	TargetPort int    `json:"targetport"`
	Protocol   string `json:"protocol"`
}
