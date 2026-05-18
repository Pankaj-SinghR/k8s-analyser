package main

import (
	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	"k8s.io/client-go/kubernetes"
)

type Scanner struct {
	Rules []rules.Rule
}

func (s *Scanner) Scan(client *kubernetes.Clientset) ([]rules.Finding, error) {
	var findings []rules.Finding
	for _, rule := range s.Rules {
		ruleFindings, err := rule.Check(client)
		if err != nil {
			return nil, err
		}
		findings = append(findings, ruleFindings...)
	}
	return findings, nil
}
