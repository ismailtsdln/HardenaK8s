# HardenaK8s 🛡️

HardenaK8s is a robust Kubernetes security auditing and hardening CLI tool written in Go. It helps you identify security misconfigurations in your clusters and provides actionable recommendations to harden your environment.

## Features
- **Comprehensive Scanning**: Audit Pods, RBAC, NetworkPolicies, and more.
- **CIS Benchmarks**: Predefined rules based on industry-standard security benchmarks.
- **Modular Policy Engine**: Support for custom YAML-based policy definitions.
- **Structured Output**: Generate reports in JSON, YAML, and HTML formats.
- **Actionable Remediations**: The `fix` command suggests or applies security improvements.

## Installation

### Prerequisites
- Go 1.21 or later
- Access to a Kubernetes cluster (kubeconfig or in-cluster)

### Install via Go
```bash
go install github.com/ismailtsdln/HardenaK8s@latest
```

### Build from source
```bash
git clone https://github.com/ismailtsdln/HardenaK8s.git
cd HardenaK8s
go build -o hardena
```

## Usage

### Scan the cluster
```bash
./hardena scan --output json
```

### Generate a report from previous results
```bash
./hardena report --input scan-results.json --output yaml
```

### Apply security fixes (Dry Run)
```bash
./hardena fix --dry-run
```

## Command Reference

| Command | Description | Flags |
|---------|-------------|-------|
| `scan`  | Scans the cluster | `--namespace`, `--all-namespaces`, `-o` |
| `report`| Generates a report | `--input`, `--output-dir`, `-o` |
| `fix`   | Applies fixes | `--dry-run` |

## CI/CD Integration
HardenaK8s can be easily integrated into your CI/CD pipelines to ensure continuous security auditing.

## License
MIT License.
