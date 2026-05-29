package k8s

import (
	"context"
	"fmt"

	"github.com/Pankaj-SinghR/k8s-analyser/rules"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type HostPathVolume struct {
	rules.RuleInfo
}

func (h HostPathVolume) Info() rules.RuleInfo {
	return h.RuleInfo
}

func NewHostPathVolume(rule rules.RuleInfo) HostPathVolume {
	return HostPathVolume{
		RuleInfo: rule,
	}
}

func (h HostPathVolume) Recommendation() string {
	return `Avoid using HostPath volumes unless absolutely necessary. HostPath volumes allow containers to mount files or directories
	directly from the host node filesystem, which can expose sensitive host resources and increase the risk of container escape or node compromise.

Recommendation:
- Use other volume types like PersistentVolumeClaims (PVCs) or emptyDir instead of HostPath
- If HostPath is necessary, ensure that the mounted path is read-only and does not contain sensitive data

Example:
spec:
  volumes:
  - name: my-volume
	emptyDir: {}
`
}

func (h HostPathVolume) Check(client *kubernetes.Clientset) ([]rules.Finding, error) {
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
			for _, volume := range pod.Spec.Volumes {
				if volume.HostPath != nil {
					finding := rules.Finding{
						ID:          h.ID,
						Resource:    fmt.Sprintf("%s in Namespace %s", pod.Name, ns.Name),
						Description: h.Description,
						Severity:    h.Severity,
					}
					findings = append(findings, finding)
				}
			}
		}
	}

	return findings, nil
}
