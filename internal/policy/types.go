package policy

import (
	"context"

	"github.com/ismailtsdln/HardenaK8s/internal/k8s"
)

// Severity represents the severity of a security issue
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
)

// Issue represents a security finding
type Issue struct {
	ID          string   `json:"id" yaml:"id"`
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description" yaml:"description"`
	Severity    Severity `json:"severity" yaml:"severity"`
	Resource    string   `json:"resource" yaml:"resource"`
	Namespace   string   `json:"namespace" yaml:"namespace"`
	Remediation string   `json:"remediation" yaml:"remediation"`
	Category    string   `json:"category" yaml:"category"`
}

// Result contains the outcome of a scan
type Result struct {
	Issues []Issue `json:"issues" yaml:"issues"`
	Stats  Stats   `json:"stats" yaml:"stats"`
}

// Stats holds summary statistics of the scan
type Stats struct {
	TotalIssues      int              `json:"total_issues" yaml:"total_issues"`
	SeverityCount    map[Severity]int `json:"severity_count" yaml:"severity_count"`
	ResourcesScanned int              `json:"resources_scanned" yaml:"resources_scanned"`
}

// Scanner defines the interface for resource-specific scanners
type Scanner interface {
	Scan(ctx context.Context, client *k8s.Client, namespace string) ([]Issue, error)
}
