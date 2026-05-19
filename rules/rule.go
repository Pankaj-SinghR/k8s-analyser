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
	Info() RuleInfo
	Recommendation() string
	Check(*kubernetes.Clientset) ([]Finding, error)
}

type RuleInfo struct {
	ID             string   `json:"rule_id"`
	Description    string   `json:"description"`
	Severity       Severity `json:"severity"`
	Recommendation string   `json:"recommendation"`
	Name           string   `json:"rule_name"`
}
