package k8s

import (
	"context"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HostNetworkEnable struct {
	rules.RuleInfo
}

func NewHostNetworkEnable(info rules.RuleInfo) HostNetworkEnable {
	return HostNetworkEnable{
		RuleInfo: info,
	}
}

func (r HostNetworkEnable) Info() rules.RuleInfo {
	return r.RuleInfo
}

func (r HostNetworkEnable) Recommendation() string {
	return `Avoid using hostNetwork: true unless absolutely necessary.
 Using the host network namespace allows pods to access the node's network interfaces directly,
 which can increase the attack surface and bypass network isolation.

 Recommendation:
 - Remove hostNetwork: true from the pod spec
 - Use Kubernetes Services for communication
 - Use NetworkPolicies for controlled traffic access

 Example:
 spec:
   hostNetwork: false
`
}

func (r HostNetworkEnable) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
			if pod.Spec.HostNetwork {
				findings = append(findings, rules.Finding{
					ID:          r.ID,
					Description: r.Description,
					Severity:    r.Severity,
					Resource:    pod.Name + " in namespace " + ns.Name,
				})
			}
		}
	}

	return findings, nil
}
