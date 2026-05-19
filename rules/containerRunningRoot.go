package rules

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ContainerRunningRoot struct {
}

func (r ContainerRunningRoot) Name() string {
	return "Container Running as Root"
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
						ID:          "CKV_K8S_2",
						Description: "Running containers as root can pose security risks. It is recommended to run containers with a non-root user.",
						Severity:    "HIGH",
						Resource:    pod.Name + " in namespace " + ns.Name,
					}
					findings = append(findings, found)
				}
			}
		}
	}

	return findings, nil
}
