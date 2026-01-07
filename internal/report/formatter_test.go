package report

import (
	"testing"

	"github.com/ismailtsdln/HardenaK8s/internal/policy"
)

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{}
	result := &policy.Result{
		Stats: policy.Stats{TotalIssues: 0},
	}

	data, err := formatter.Format(result)
	if err != nil {
		t.Fatalf("failed to format JSON: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty JSON data")
	}
}

func TestGetFormatter(t *testing.T) {
	_, err := GetFormatter("json")
	if err != nil {
		t.Errorf("expected JSON formatter, got error: %v", err)
	}

	_, err = GetFormatter("invalid")
	if err == nil {
		t.Error("expected error for invalid formatter, got nil")
	}
}
