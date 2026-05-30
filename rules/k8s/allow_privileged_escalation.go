package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type AllowPrivilegeEscalation struct {
	rules.RuleInfo
}

func (a AllowPrivilegeEscalation) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewAllowPrivilegeEscalation(rule rules.RuleInfo) AllowPrivilegeEscalation {
	return AllowPrivilegeEscalation{
		RuleInfo: rule,
	}
}

func (a AllowPrivilegeEscalation) Recommendation() string {
	return `Disable privilege escalation in the container security context.
 Example:
 securityContext:
   allowPrivilegeEscalation: false

 Additionally, apply least privilege practices by:
 - Running containers as non-root
 - Dropping unnecessary Linux capabilities
 - Avoiding privileged containers
`
}

func (a AllowPrivilegeEscalation) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
				if container.SecurityContext != nil && container.SecurityContext.AllowPrivilegeEscalation != nil && *container.SecurityContext.AllowPrivilegeEscalation {
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
