package rules

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ContainerRunningRoot struct {
	RuleInfo
}

func (c ContainerRunningRoot) Info() RuleInfo {
	return c.RuleInfo
}

func NewContainerRunningRoot(rule RuleInfo) ContainerRunningRoot {
	return ContainerRunningRoot{
		RuleInfo: rule,
	}
}

func (c ContainerRunningRoot) Recommendation() string {
	//   Set:
	// securityContext:
	//   runAsNonRoot: true
	rr := fmt.Sprintf(`Set the securityContext for the container to run as a non-root user.
 Example:
  securityContext:
   runAsNonRoot: true`)
	return rr
}

func (r ContainerRunningRoot) Check(client *kubernetes.Clientset) ([]Finding, error) {
	// check all namespaces, then check all pods in each namespace, and check if any container is running as root
	// securityContext.RunAsUser is 0 or nil (default is 0) means running as root

	namespace, err := client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var findings []Finding

	for _, ns := range namespace.Items {
		pods, err := client.CoreV1().Pods(ns.Name).List(context.Background(), v1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for _, pod := range pods.Items {
			for _, container := range pod.Spec.Containers {
				if container.SecurityContext != nil && (container.SecurityContext.RunAsUser == nil || *container.SecurityContext.RunAsUser == 0) {
					found := Finding{
						ID:          r.ID,
						Description: r.Description,
						Severity:    r.Severity,
						Resource:    pod.Name + " in namespace " + ns.Name,
					}
					findings = append(findings, found)
				}
			}
		}
	}

	return findings, nil
}
