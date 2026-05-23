package main

import (
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	"k8s.io/client-go/kubernetes"
)

type Scanner struct {
	Rules []rules.Rule
}

func (s *Scanner) Scan(client *kubernetes.Clientset) ([]rules.Finding, error) {
	var findings []rules.Finding
	fmt.Println("Findings")
	for _, rule := range s.Rules {
		ruleFindings, err := rule.Check(client)
		if err != nil {
			return nil, err
		}
		// if len(ruleFindings) == 0 {
		// 	continue
		// }
		info := rule.Info()
		fmt.Println("════════════════════════════════════════")
		fmt.Printf("\n")
		fmt.Printf("[%s] %s\n", info.Severity, info.ID)
		fmt.Printf("────────────────────────────────────────\n")
		fmt.Printf("Rule		: %s\n", info.Name)
		fmt.Printf("Description	: %s\n", info.Description)
		fmt.Printf("\n")
		fmt.Println("Affected Resources:")
		for _, finding := range ruleFindings {
			fmt.Printf(" • %s\n", finding.Resource)
		}
		fmt.Printf("\n")
		fmt.Println("Recommendation	:")
		fmt.Printf(" %s\n", rule.Recommendation())
		fmt.Printf("\n")
		findings = append(findings, ruleFindings...)
	}
	return findings, nil
}
