package rules

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CheckLatestTag struct{}

func (c CheckLatestTag) ID() string {
	return "CKV_K8S_1"
}

func (c CheckLatestTag) Description() string {
	return "Using 'latest' tag in container image is not recommended as it can lead to unpredictable deployments."
}

func (c CheckLatestTag) Severity() Severity {
	return High
}

// recommendation is to use how to fix the issue, for example, if the rule is about using 'latest' tag,
// then the recommendation would be to use a specific tag instead of 'latest'
func (c CheckLatestTag) Recommendation() string {
	rr := fmt.Sprintf(`Use a fixed image tag instead of 'latest'
 Example: nginx:1.27.1`)
	return rr
}

func (c CheckLatestTag) Name() string {
	return "Check for latest tag usage"
}

func (c CheckLatestTag) Check(client *kubernetes.Clientset) ([]Finding, error) {
	// get all namespaces, then get all pods in each namespace, and check if any container is using 'latest' tag
	// use pagination to get all namespaces and pods if there are many
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
				if container.Image != "" && strings.Split(container.Image, ":")[1] == "latest" {
					found := Finding{
						ID:          c.ID(),
						Description: c.Description(),
						Severity:    c.Severity(),
						Resource:    pod.Name + " in namespace " + ns.Name,
					}
					findings = append(findings, found)
				}
			}
		}
	}

	return findings, nil
}
