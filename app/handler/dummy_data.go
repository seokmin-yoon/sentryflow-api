package handler

import (
	
	"sentryflow-api/app/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
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
			Pods: []model.Pod{
				{
					TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx-pod-1",
						Namespace: "default",
						Labels:    map[string]string{"app": "nginx"},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "nginx-container",
								Image: "nginx:latest",
								Ports: []v1.ContainerPort{
									{ContainerPort: 80},
								},
							},
						},
					},
					Status: v1.PodStatus{
						Phase:  v1.PodRunning,
						PodIP:  "10.244.1.10",
						HostIP: "192.168.1.100",
					},
				},
				{
					TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "api-server",
						Namespace: "dev",
						Labels:    map[string]string{"app": "api"},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "api-container",
								Image: "my-api:v1.0",
								Ports: []v1.ContainerPort{
									{ContainerPort: 8080},
								},
							},
						},
					},
					Status: v1.PodStatus{
						Phase:  v1.PodRunning,
						PodIP:  "10.244.2.15",
						HostIP: "192.168.1.101",
					},
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
					TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "database-pod",
						Namespace: "staging",
						Labels:    map[string]string{"app": "database"},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "db-container",
								Image: "postgres:14",
								Ports: []v1.ContainerPort{
									{ContainerPort: 5432},
								},
							},
						},
					},
					Status: v1.PodStatus{
						Phase:  v1.PodPending,
						PodIP:  "10.244.3.20",
						HostIP: "192.168.1.102",
					},
				},
			},
		},
	}
)
