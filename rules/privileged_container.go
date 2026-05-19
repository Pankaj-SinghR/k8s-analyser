package rules

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CheckPrivilegedContainer struct {
	RuleInfo
}

func (c CheckPrivilegedContainer) Info() RuleInfo {
	return c.RuleInfo
}

func NewCheckPrivilegedContainer(info RuleInfo) CheckPrivilegedContainer {
	return CheckPrivilegedContainer{
		RuleInfo: info,
	}
}

// recommendation is to use how to fix the issue, for example, if the rule is about using 'latest' tag,
// then the recommendation would be to use a specific tag instead of 'latest'
func (c CheckPrivilegedContainer) Recommendation() string {
	rr := fmt.Sprintf(`Avoid running containers with privileged access. Instead, use specific capabilities or security contexts to limit the permissions of the container.
 Example:
 apiVersion: v1
 kind: Pod
 metadata:
   name: example-pod
 spec:
   containers:
   - name: example-container
	 image: nginx:1.27.1
	 securityContext:
	   privileged: false`)
	return rr
}

func (c CheckPrivilegedContainer) Check(client *kubernetes.Clientset) ([]Finding, error) {
	namespace, err := client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	var findings []Finding

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
				if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged == true {
					found := Finding{
						ID:          c.ID,
						Description: c.Description,
						Severity:    c.Severity,
						Resource:    pod.Name + " in namespace " + ns.Name,
					}
					findings = append(findings, found)
				}
			}
		}
	}

	return findings, nil
}
