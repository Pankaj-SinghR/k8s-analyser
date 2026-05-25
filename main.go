package main

import (
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

	// create scanner with rules
	scanner := Scanner{
		Rules: []rules.Rule{
			rules.CheckLatestTag{},
			rules.ContainerRunningRoot{},
			rules.CheckPrivilegedContainer{},
		},
	}

	// run scanner
	_, err = scanner.Scan(client)

	if err != nil {
		println("Error running scanner:", err.Error())
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
