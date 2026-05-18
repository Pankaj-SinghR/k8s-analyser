package rules

import "k8s.io/client-go/kubernetes"

type Finding struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Resource    string `json:"resource"`
}

type Rule interface {
	Name() string
	Check(*kubernetes.Clientset) ([]Finding, error)
}
