package policy

import (
	"context"
	"fmt"

	"github.com/ismailtsdln/HardenaK8s/internal/k8s"
	"github.com/ismailtsdln/HardenaK8s/internal/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Engine coordinates the scanning process
type Engine struct {
	client   *k8s.Client
	scanners []Scanner
}

// NewEngine creates a new policy engine
func NewEngine(client *k8s.Client) *Engine {
	return &Engine{
		client: client,
		scanners: []Scanner{
			&PodScanner{}, // Add more scanners here
		},
	}
}

// Run executes all registered scanners
func (e *Engine) Run(ctx context.Context, namespace string) (*Result, error) {
	result := &Result{
		Issues: []Issue{},
		Stats: Stats{
			SeverityCount: map[Severity]int{
				SeverityCritical: 0,
				SeverityHigh:     0,
				SeverityMedium:   0,
				SeverityLow:      0,
				SeverityInfo:     0,
			},
		},
	}

	for _, scanner := range e.scanners {
		issues, err := scanner.Scan(ctx, e.client, namespace)
		if err != nil {
			logger.Error("Scanner failed", "error", err)
			continue
		}

		for _, issue := range issues {
			result.Issues = append(result.Issues, issue)
			result.Stats.TotalIssues++
			result.Stats.SeverityCount[issue.Severity]++
		}
	}

	return result, nil
}

// PodScanner audits Pod configurations
type PodScanner struct{}

// Scan audits pods in the given namespace
func (p *PodScanner) Scan(ctx context.Context, client *k8s.Client, namespace string) ([]Issue, error) {
	var issues []Issue

	pods, err := client.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	for _, pod := range pods.Items {
		// Example check: Privileged container
		for _, container := range pod.Spec.Containers {
			if container.SecurityContext != nil && container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
				issues = append(issues, Issue{
					ID:          "HK-001",
					Title:       "Privileged Container Detected",
					Description: fmt.Sprintf("Pod %s in namespace %s has a privileged container: %s", pod.Name, pod.Namespace, container.Name),
					Severity:    SeverityCritical,
					Resource:    pod.Name,
					Namespace:   pod.Namespace,
					Remediation: "Remove 'privileged: true' from securityContext.",
					Category:    "Pod Security",
				})
			}
		}
	}

	return issues, nil
}
