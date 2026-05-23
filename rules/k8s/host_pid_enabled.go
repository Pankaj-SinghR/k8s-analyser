package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HostPIDEnabled struct {
	rules.RuleInfo
}

func (h HostPIDEnabled) Info() rules.RuleInfo {
	return h.RuleInfo
}

func NewHostPIDEnabled(rule rules.RuleInfo) HostPIDEnabled {
	return HostPIDEnabled{
		RuleInfo: rule,
	}
}

func (h HostPIDEnabled) Recommendation() string {
	return `Avoid using hostPID: true unless absolutely necessary.
Using the host PID namespace allows containers to view and interact
with processes running on the host node, which can lead to privilege
escalation and container breakout risks.

Recommendation:
- Remove hostPID: true from the pod spec
- Use isolated namespaces
- Follow the principle of least privilege

Example:
spec:
  hostPID: false
`
}

func (h HostPIDEnabled) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
	// get all namespaces, then get all pods in each namespace, and check if any container is using 'latest' tag
	// use pagination to get all namespaces and pods if there are many
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
			if pod.Spec.HostPID {
				finding := rules.Finding{
					ID:          h.ID,
					Description: h.Description,
					Severity:    h.Severity,
					Resource:    fmt.Sprintf("%s in Namespace %s", pod.Name, ns.Name),
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings, nil
}
