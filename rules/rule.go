package rules

import "k8s.io/client-go/kubernetes"

type Severity string

const (
	Low    Severity = "LOW"
	Medium Severity = "MEDIUM"
	High   Severity = "HIGH"
)

type Finding struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
	Resource    string   `json:"resource"`
}

type Rule interface {
	ID() string
	Description() string
	Severity() Severity
	Recommendation() string
	Name() string
	Check(*kubernetes.Clientset) ([]Finding, error)
}
