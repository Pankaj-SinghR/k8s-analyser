package main

import "k8s.io/client-go/kubernetes"

type Scanner struct {
	Rules []Rule
}

func (s *Scanner) Scan(client *kubernetes.Clientset) []Finding {
	var findings []Finding
	for _, rule := range s.Rules {
		ruleFindings := rule.Check(client)
		findings = append(findings, ruleFindings...)
	}
	return findings
}
