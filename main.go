package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// load config, from ~/.kube/config
	// create clientset
	client, err := NewClientset()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully created Kubernetes clientset")

	// create scanner with rules
	scanner := Scanner{
		Rules: []rules.Rule{
			rules.CheckLatestTag{},
			rules.ContainerRunningRoot{},
		},
	}

	// run scanner
	findings, err := scanner.Scan(client)

	if err != nil {
		println("Error running scanner:", err.Error())
	}

	for _, finding := range findings {
		fmt.Printf("ID: %s, Description: %s, Severity: %s, Resource: %s \n", finding.ID, finding.Description, finding.Severity, finding.Resource)
	}
}

func NewClientset() (*kubernetes.Clientset, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
