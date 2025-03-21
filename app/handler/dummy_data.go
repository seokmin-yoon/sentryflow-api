package handler

import (
	"time"

	"sentryflow-api/app/model"
)

// stored data example
var (
	users = []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}

	clusters = []model.Cluster{
		{
			Name: "Cluster-A",
			Namespaces: []model.Namespace{
				{Name: "default", Phase: "Active"},
				{Name: "kube-system", Phase: "Active"},
				{Name: "dev", Phase: "Active"},
				{Name: "prod", Phase: "Active"},
			},
			Pods: []model.Pod{
				{
					ID:          "0",
					Name:        "sleep-7656cf8794-tpc4k",
					PodIP:       "10.244.0.12",
					ClusterName: "Cluster-A",
					Namespace:   "default",
					QosClass:    "BestEffort",
					Status:      "Running",
					Timestamp:   time.Now().Format(time.RFC3339),
				},
				{
					ID:          "1",
					Name:        "httpbin-6c4d945c8d-p2tnw",
					PodIP:       "10.244.0.11",
					ClusterName: "Cluster-A",
					Namespace:   "default",
					QosClass:    "BestEffort",
					Status:      "Running",
					Timestamp:   time.Now().Format(time.RFC3339),
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
			Pods: []model.Pod{
				{
					ID:          "0",
					Name:        "sentryflow-api-bc84df6ff-kg8cl",
					PodIP:       "10.244.0.13",
					ClusterName: "Cluster-B",
					Namespace:   "sentryflow",
					QosClass:    "BestEffort",
					Status:      "Running",
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
		},
	}
)
