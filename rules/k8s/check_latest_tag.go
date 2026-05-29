package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CheckLatestTag struct {
	rules.RuleInfo
}

func (c CheckLatestTag) Info() rules.RuleInfo {
	return c.RuleInfo
}

func NewCheckLatestTag(rule rules.RuleInfo) CheckLatestTag {
	return CheckLatestTag{
		RuleInfo: rule,
	}
}

// recommendation is to use how to fix the issue, for example, if the rule is about using 'latest' tag,
// then the recommendation would be to use a specific tag instead of 'latest'
func (c CheckLatestTag) Recommendation() string {
	rr := fmt.Sprintf(`Use a fixed image tag instead of 'latest'
 Example: nginx:1.27.1`)
	return rr
}

func (c CheckLatestTag) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
			for _, container := range pod.Spec.Containers {
				if container.Image != "" && strings.Split(container.Image, ":")[1] == "latest" {
					found := rules.Finding{
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
