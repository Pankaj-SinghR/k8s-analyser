package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type AutoMountServiceToken struct {
	rules.RuleInfo
}

func (a AutoMountServiceToken) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewAutoMountServiceToken(rule rules.RuleInfo) AutoMountServiceToken {
	return AutoMountServiceToken{
		RuleInfo: rule,
	}
}

func (a AutoMountServiceToken) Recommendation() string {
	recommendation := `Avoid using automountServiceAccountToken: true unless absolutely necessary. When automountServiceAccountToken 
	is set to true, Kubernetes automatically mounts a service account token into the pod, which can be used by attackers to access the 
	Kubernetes API and potentially escalate privileges.
Recommendation:
- Set automountServiceAccountToken: false in the pod spec
- If the service account token is required, ensure that it has minimal permissions and follows the principle of least privilege

Example:
spec:
  automountServiceAccountToken: false
`
	return recommendation
}

func (a AutoMountServiceToken) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
			if pod.Spec.AutomountServiceAccountToken != nil && *pod.Spec.AutomountServiceAccountToken {
				finding := rules.Finding{
					ID:          a.ID,
					Description: a.Description,
					Severity:    a.Severity,
					Resource:    fmt.Sprintf("%s in Namespace %s", pod.Name, ns.Name),
				}
				findings = append(findings, finding)
			}
		}
	}

	return findings, nil
}
