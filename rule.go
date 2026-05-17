package main

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Finding struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Resource    string `json:"resource"`
}

type Rule interface {
	Name() string
	Check(*kubernetes.Clientset) ([]Finding, error)
}

type CheckLatestTag struct{}

func (c CheckLatestTag) Name() string {
	return "Check for latest tag usage"
}

func (c CheckLatestTag) Check(client *kubernetes.Clientset) ([]Finding, error) {
	// get all namespaces, then get all pods in each namespace, and check if any container is using 'latest' tag
	// use pagination to get all namespaces and pods if there are many
	namespace, err := client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})

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
					return []Finding{
						{
							ID:          "CKV_K8S_1",
							Description: "Using 'latest' tag in container image is not recommended as it can lead to unpredictable deployments.",
							Severity:    "HIGH",
							Resource:    pod.Name + " in namespace " + ns.Name,
						},
					}, nil
				}
			}
		}
	}

	return nil, nil
}
