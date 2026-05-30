package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MissingCPULimits struct {
	rules.RuleInfo
}

func (a MissingCPULimits) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewMissingCPULimits(rule rules.RuleInfo) MissingCPULimits {
	return MissingCPULimits{
		RuleInfo: rule,
	}
}

func (a MissingCPULimits) Recommendation() string {
	recommendation := `Containers without CPU limits can consume excessive CPU resources, potentially affecting cluster stability and other workloads. It is recommended to define CPU limits for all containers.

Recommendation:
- Configure CPU limits for all containers
- Use appropriate CPU requests and limits based on application requirements

Example:
spec:
  containers:
  - name: example-container
    image: nginx:1.27.1
    resources:
      limits:
        cpu: "500m"
      requests:
        cpu: "100m"
	`
	return recommendation
}

func (a MissingCPULimits) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
				if container.Resources.Limits.Cpu().IsZero() {
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
