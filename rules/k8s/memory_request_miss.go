package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MissingMemoryRequests struct {
	rules.RuleInfo
}

func (a MissingMemoryRequests) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewMissingMemoryRequests(rule rules.RuleInfo) MissingMemoryRequests {
	return MissingMemoryRequests{
		RuleInfo: rule,
	}
}

func (a MissingMemoryRequests) Recommendation() string {
	recommendation := `Containers without memory requests may lead to inefficient resource scheduling and increase the risk of resource contention. It is recommended to define memory requests for all containers to ensure predictable workload performance.

Recommendation:
- Configure memory requests for all containers
- Use appropriate memory requests based on application requirements

Example:
spec:
  containers:
  - name: example-container
    image: nginx:1.27.1
    resources:
      requests:
        memory: "128Mi"
	`
	return recommendation
}

func (a MissingMemoryRequests) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
				if container.Resources.Requests.Memory().IsZero() {
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
