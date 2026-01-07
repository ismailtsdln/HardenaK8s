package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"gopkg.in/yaml.v3"
)

// Formatter defines the interface for report formatters
type Formatter interface {
	Format(result *policy.Result) ([]byte, error)
}

// JSONFormatter implements Formatter for JSON output
type JSONFormatter struct{}

func (f *JSONFormatter) Format(result *policy.Result) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

// YAMLFormatter implements Formatter for YAML output
type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(result *policy.Result) ([]byte, error) {
	return yaml.Marshal(result)
}

// TextFormatter is a fallback for CLI output
type TextFormatter struct{}

func (f *TextFormatter) Format(result *policy.Result) ([]byte, error) {
	// This is typically handled by the CLI logic directly for vibrancy
	return []byte("Summary: " + fmt.Sprintf("%d issues found", result.Stats.TotalIssues)), nil
}

// SaveToFile writes the formatted report to a file
func SaveToFile(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0644)
}

// GetFormatter returns the appropriate formatter based on the format string
func GetFormatter(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSONFormatter{}, nil
	case "yaml":
		return &YAMLFormatter{}, nil
	case "html":
		return &HTMLFormatter{}, nil
	case "text":
		return &TextFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
