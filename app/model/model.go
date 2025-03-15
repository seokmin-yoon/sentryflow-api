package model

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APILog struct {
	ID          uint64            `json:"id"`
	TimeStamp   string            `json:"timeStamp"`

	SrcNamespace string            `json:"srcNamespace"`
	SrcName      string            `json:"srcName"`
	SrcLabel     map[string]string `json:"srcLabel"`

	SrcType string `json:"srcType"`
	SrcIP   string `json:"srcIP"`
	SrcPort string `json:"srcPort"`

	DstNamespace string            `json:"dstNamespace"`
	DstName      string            `json:"dstName"`
	DstLabel     map[string]string `json:"dstLabel"`

	DstType string `json:"dstType"`
	DstIP   string `json:"dstIP"`
	DstPort string `json:"dstPort"`

	Protocol     string `json:"protocol"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	ResponseCode int32  `json:"responseCode"`
}