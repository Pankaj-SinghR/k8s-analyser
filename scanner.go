package main

import "k8s.io/client-go/kubernetes"

type Scanner struct {
	Rules []Rule
}

func (s *Scanner) Scan(client *kubernetes.Clientset) ([]Finding, error) {
	var findings []Finding
	for _, rule := range s.Rules {
		ruleFindings, err := rule.Check(client)
		if err != nil {
			return nil, err
		}
		findings = append(findings, ruleFindings...)
	}
	return findings, nil
}
