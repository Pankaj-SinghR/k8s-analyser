package main

import (
	"encoding/json"
	"os"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
)

func init() {
	println("Starting Kubernetes Analyser...")
}

func MapRuleWithJSON() map[string]rules.Rule {
	file, err := os.Open("rules.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var rulesMap map[string]rules.RuleInfo

	err = json.NewDecoder(file).Decode(&rulesMap)
	if err != nil {
		panic(err)
	}

	return map[string]rules.Rule{
		"CKV_K8S_1": rules.NewCheckLatestTag(rulesMap["CKV_K8S_1"]),
		"CKV_K8S_2": rules.NewContainerRunningRoot(rulesMap["CKV_K8S_2"]),
		"CKV_K8S_3": rules.NewCheckPrivilegedContainer(rulesMap["CKV_K8S_3"]),
		"CKV_K8S_4": rules.NewHostNetworkEnable(rulesMap["CKV_K8S_4"]),
	}
}
