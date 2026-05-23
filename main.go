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
			MapRuleWithJSON()["CKV_K8S_1"],
			MapRuleWithJSON()["CKV_K8S_2"],
			MapRuleWithJSON()["CKV_K8S_3"],
			MapRuleWithJSON()["CKV_K8S_4"],
			MapRuleWithJSON()["CKV_K8S_5"],
			MapRuleWithJSON()["CKV_K8S_6"],
			MapRuleWithJSON()["CKV_K8S_7"],
			MapRuleWithJSON()["CKV_K8S_8"],
			MapRuleWithJSON()["CKV_K8S_9"],
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
