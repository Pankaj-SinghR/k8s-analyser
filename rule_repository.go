package main

import (
	"encoding/json"
	"os"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	"github.com/Pankaj-SinghR/k8s-analyser/rules/k8s"
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
		"CKV_K8S_1": k8s.NewCheckLatestTag(rulesMap["CKV_K8S_1"]),
		"CKV_K8S_2": k8s.NewContainerRunningRoot(rulesMap["CKV_K8S_2"]),
		"CKV_K8S_3": k8s.NewCheckPrivilegedContainer(rulesMap["CKV_K8S_3"]),
		"CKV_K8S_4": k8s.NewHostNetworkEnable(rulesMap["CKV_K8S_4"]),
		"CKV_K8S_5": k8s.NewHostPIDEnabled(rulesMap["CKV_K8S_5"]),
		"CKV_K8S_6": k8s.NewHostIPCEnabled(rulesMap["CKV_K8S_6"]),
		"CKV_K8S_7": k8s.NewAllowPrivilegeEscalation(rulesMap["CKV_K8S_7"]),
	}
}
