# k8s-misconfig-finder

A lightweight Kubernetes misconfiguration detection CLI for identifying dangerous security, reliability, and operational issues in live Kubernetes clusters.

The tool connects directly to the cluster using the current kubeconfig context and scans workloads such as Pods, Deployments, and StatefulSets for common infrastructure and security anti-patterns.

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
k8s-scan cluster
```

### Example Output

```text
[HIGH] payments-api
- Container running as root
- Privileged mode enabled

[MEDIUM] auth-service
- Missing memory limits
- Missing readiness probe

[LOW] nginx
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
- Cobra CLI
- Kubernetes client-go
- Kubernetes API

---

## Roadmap

### MVP

- [x] Connect to Kubernetes cluster
- [ ] Scan Pods and Deployments
- [ ] Detect common misconfigurations
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

Kubernetes clusters often contain dangerous or inefficient configurations that are difficult to identify manually.

This project focuses on providing a simple, developer-friendly auditing tool to help platform and DevOps engineers quickly detect critical infrastructure misconfigurations in running environments.

---

## License

MIT
