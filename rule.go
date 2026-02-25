package main

import "k8s.io/client-go/kubernetes"

type Finding struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Resource    []string `json:"resource"`
}

type Rule interface {
	Name() string
	Check(*kubernetes.Clientset) []Finding
}

type CheckLatestTag struct{}

func (c CheckLatestTag) Name() string {
	return "Check for latest tag usage"
}

func (c CheckLatestTag) Check(client *kubernetes.Clientset) []Finding {
	// check logic goes here, for now we return a dummy finding
	return []Finding{
		{
			ID:          "CKV_K8S_1",
			Description: "Using 'latest' tag in container image is not recommended as it can lead to unpredictable deployments.",
			Severity:    "HIGH",
		},
	}
}
