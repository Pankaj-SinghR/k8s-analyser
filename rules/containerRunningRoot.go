package rules

import "k8s.io/client-go/kubernetes"

type ContainerRunningRoot struct {
}

func (r *ContainerRunningRoot) Name() string {
	return "Container Running as Root"
}

func (r *ContainerRunningRoot) Check(client *kubernetes.Clientset) ([]Finding, error) {
	return []Finding{
		{
			ID:          "CKV_K8S_2",
			Description: "Running containers as root can pose security risks. It is recommended to run containers with a non-root user.",
			Severity:    "HIGH",
			Resource:    "Pod/Deployment/StatefulSet/DaemonSet",
		},
	}, nil
}
