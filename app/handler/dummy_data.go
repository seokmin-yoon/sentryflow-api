package handler

import (
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
			Name:       "Cluster1",
			Namespaces: []string{"default", "kube-system", "app-space"},
		},
		{
			Name:       "Cluster2",
			Namespaces: []string{"dev", "staging", "monitoring"},
		},
	}
)
