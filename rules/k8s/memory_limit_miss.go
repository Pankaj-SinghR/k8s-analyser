package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MissingMemoryLimits struct {
	rules.RuleInfo
}

func (a MissingMemoryLimits) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewMissingMemoryLimits(rule rules.RuleInfo) MissingMemoryLimits {
	return MissingMemoryLimits{
		RuleInfo: rule,
	}
}

func (a MissingMemoryLimits) Recommendation() string {
	recommendation := `Containers without memory limits can consume excessive memory resources, potentially causing node instability or OOM (Out Of Memory) issues. It is recommended to define memory limits for all containers.

Recommendation:
- Configure memory limits for all containers
- Use appropriate memory requests and limits based on application requirements

Example:
spec:
  containers:
  - name: example-container
    image: nginx:1.27.1
    resources:
      limits:
        memory: "512Mi"
      requests:
        memory: "128Mi"
	`
	return recommendation
}

func (a MissingMemoryLimits) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
	namespace, err := client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	var findings []rules.Finding

	if err != nil {
		return nil, err
	}

	for _, ns := range namespace.Items {
		pods, err := client.CoreV1().Pods(ns.Name).List(context.Background(), v1.ListOptions{})

		if err != nil {
			return nil, err
		}

		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if container.Resources.Limits.Memory().IsZero() {
					finding := rules.Finding{
						Resource:    fmt.Sprintf("%s in Namespace %s", pod.Name, ns.Name),
						Description: a.Description,
						ID:          a.ID,
						Severity:    a.Severity,
					}
					findings = append(findings, finding)
				}
			}
		}
	}

	return findings, nil
}
