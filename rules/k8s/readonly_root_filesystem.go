package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ReadOnlyRootFileSystem struct {
	rules.RuleInfo
}

func (a ReadOnlyRootFileSystem) Info() rules.RuleInfo {
	return a.RuleInfo
}

func NewRootOnlyRootFileSystem(rule rules.RuleInfo) ReadOnlyRootFileSystem {
	return ReadOnlyRootFileSystem{
		RuleInfo: rule,
	}
}

func (a ReadOnlyRootFileSystem) Recommendation() string {
	recommendation := `Running containers with a read-only root filesystem can enhance security by preventing unauthorized modifications to the container's filesystem. It is recommended to set the root filesystem to read-only unless the application requires write access.
Recommendation:
- Set readOnlyRootFilesystem: true in the container's security context
- If the application requires write access, consider using an emptyDir volume or a PersistentVolumeClaim (PVC) for writable storage instead of allowing write access to the root filesystem

Example:
spec:
  containers:
  - name: example-container
    image: nginx:1.27.1
    securityContext:
    readOnlyRootFilesystem: true
	`
	return recommendation
}

func (a ReadOnlyRootFileSystem) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
				if container.SecurityContext != nil && (container.SecurityContext.ReadOnlyRootFilesystem == nil || !*container.SecurityContext.ReadOnlyRootFilesystem) {
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
