package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// list all namespaces
	val, _ := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	log.Printf("Namespaces: %v", val.Items)

	for _, ns := range val.Items {
		log.Printf("Namespace: %s", ns.Name)
		log.Printf("Status: %s", &ns.Status.Phase)
		pods, _ := client.CoreV1().Pods(ns.Name).List(context.Background(), metav1.ListOptions{})
		for _, pod := range pods.Items {
			log.Printf("Pod: %s, Status: %s", pod.Name, pod.Status.Phase)
		}
	}

	// for each namespace, list all pods
	// for each pod, print its name and status

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
