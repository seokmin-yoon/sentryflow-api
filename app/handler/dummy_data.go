package handler

import (
	"sentryflow-api/app/model"
	"time"
)

// stored data example
var (
	users = []model.User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}

	clusters = []model.Cluster{
		{
			Name: "Cluster-A",
			Namespaces: []model.Namespace{
				{Name: "default", Phase: "Active"},
				{Name: "kube-system", Phase: "Active"},
				{Name: "dev", Phase: "Active"},
				{Name: "prod", Phase: "Active"},
			},
			Nodes: []model.Node{
				{
					Name:             "node-1",
					KernelVersion:    "5.15.0-79-generic",
					OSImage:          "Ubuntu 22.04 LTS",
					ContainerRuntime: "containerd://1.6.18",
					KubeletVersion:   "v1.26.3",
					KubeProxyVersion: "v1.26.3",
					PodCIDR:          "10.244.0.0/24",
					ProviderID:       "aws://us-west-2/i-1234567890abcdef0",
					SystemUUID:       "4f82cb20-5b7f-11ea-9fbb-00155d6402da",
					InternalIP:       "192.168.1.10",
					CreatedAt:        time.Now(),
				},
				{
					Name:             "node-2",
					KernelVersion:    "5.15.0-79-generic",
					OSImage:          "Ubuntu 22.04 LTS",
					ContainerRuntime: "containerd://1.6.18",
					KubeletVersion:   "v1.26.3",
					KubeProxyVersion: "v1.26.3",
					PodCIDR:          "10.244.1.0/24",
					ProviderID:       "aws://us-west-2/i-1234567890abcdef1",
					SystemUUID:       "4f82cb20-5b7f-11ea-9fbb-00155d6402db",
					InternalIP:       "192.168.1.11",
					CreatedAt:        time.Now(),
				},
			},
		},
		{
			Name: "Cluster-B",
			Namespaces: []model.Namespace{
				{Name: "default", Phase: "Active"},
				{Name: "kube-public", Phase: "Active"},
				{Name: "staging", Phase: "Active"},
			},
			Nodes: []model.Node{
				{
					Name:             "node-3",
					KernelVersion:    "5.10.0-60-generic",
					OSImage:          "Ubuntu 20.04 LTS",
					ContainerRuntime: "docker://20.10.18",
					KubeletVersion:   "v1.25.6",
					KubeProxyVersion: "v1.25.6",
					PodCIDR:          "10.245.0.0/24",
					ProviderID:       "gce://us-central1/i-234567890abcdef0",
					SystemUUID:       "5f92cb20-6b8f-12eb-8fcd-00265d6403da",
					InternalIP:       "192.168.2.10",
					CreatedAt:        time.Now(),
				},
			},
		},
	}
)
