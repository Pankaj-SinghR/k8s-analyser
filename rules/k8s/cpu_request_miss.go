package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MissingCPURequests struct {
	rules.RuleInfo
}

func (a MissingCPURequests) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewMissingCPURequests(rule rules.RuleInfo) MissingCPURequests {
	return MissingCPURequests{
		RuleInfo: rule,
	}
}

func (a MissingCPURequests) Recommendation() string {
	recommendation := `Containers without CPU requests may lead to inefficient resource scheduling and unstable workload performance. It is recommended to define CPU requests for all containers to ensure proper resource allocation.
Recommendation:
- Configure CPU requests for all containers
- Use appropriate CPU requests based on application requirements

Example:
spec:
  containers:
  - name: example-container
    image: nginx:1.27.1
    resources:
      requests:
        cpu: "100m"
	`
	return recommendation
}

func (a MissingCPURequests) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
				if container.Resources.Requests.Cpu().IsZero() {
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
