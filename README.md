# k8s-analyser

A lightweight Kubernetes security and reliability analysis CLI for identifying dangerous infrastructure risks, unsafe workload configurations, and operational anti-patterns in live Kubernetes clusters.

The tool connects directly to the cluster using the current kubeconfig context and scans workloads such as Pods, Deployments, and StatefulSets to detect security weaknesses, reliability concerns, and Kubernetes best-practice violations.

---

## Features

- Scan live Kubernetes clusters
- Detect dangerous container configurations
- Identify missing resource limits and probes
- Group findings by severity
- Fast CLI-first workflow
- No UI dependencies

---

## Misconfigurations Detected

### Security

- Containers running as root
- Privileged containers
- Host network enabled
- Host PID / IPC enabled
- AllowPrivilegeEscalation enabled
- Mutable image tags (`latest`)

### Reliability

- Missing liveness probes
- Missing readiness probes
- Missing CPU limits
- Missing memory limits
- Missing resource requests

### Best Practices

- Workloads running in the default namespace
- Missing labels
- Missing annotations

---

## Example Usage

### Scan Entire Cluster

```bash
k8s-analyser cluster
```

### Example Output

```text
[HIGH]
- Container running as root
- Privileged mode enabled

[MEDIUM]
- Missing memory limits
- Missing readiness probe

[LOW]
- Running in default namespace
```

---

## Architecture

```text
Kubernetes API
      ↓
Workload Fetcher
      ↓
Rule Engine
      ↓
Findings Analyzer
      ↓
CLI Output
```

---

## Tech Stack

- Golang
- Kubernetes client-go
- Kubernetes API

---

## Roadmap

### MVP

- [x] Connect to Kubernetes cluster
- [x] Scan Pods and Deployments
- [x] Detect common misconfigurations
- [ ] Severity-based findings
- [ ] CLI output formatting

### Future Improvements

- [ ] Namespace filtering
- [ ] JSON output
- [ ] CI/CD integration
- [ ] Slack notifications
- [ ] Custom rule engine
- [ ] Export reports

---

## Why This Project?

Kubernetes clusters often contain unsafe or inefficient configurations that are difficult to identify manually.

This project provides a lightweight CLI to help engineers quickly detect security, reliability, and operational issues in running Kubernetes environments.

---

## License

MIT
