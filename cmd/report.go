package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ismailtsdln/HardenaK8s/internal/logger"
	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"github.com/ismailtsdln/HardenaK8s/internal/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report from scan results",
	Long: `The report command processes the results of a previous scan 
and generates a report in the specified format (JSON, YAML, or HTML).`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, _ := cmd.Flags().GetString("input")
		outputDir, _ := cmd.Flags().GetString("output-dir")
		outputFormat := viper.GetString("output")

		logger.Log.Info("Generating report", "input", inputFile, "format", outputFormat)

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			logger.Log.Error("Failed to create output directory", "error", err)
			return
		}

		// Read input results
		data, err := os.ReadFile(inputFile)
		if err != nil {
			logger.Log.Error("Failed to read input file", "error", err)
			return
		}

		var result policy.Result
		if err := json.Unmarshal(data, &result); err != nil {
			logger.Log.Error("Failed to parse scan results", "error", err)
			return
		}

		// Format report
		formatter, err := report.GetFormatter(outputFormat)
		if err != nil {
			logger.Log.Warn("Invalid output format, defaulting to JSON", "format", outputFormat)
			formatter = &report.JSONFormatter{}
		}

		outputData, err := formatter.Format(&result)
		if err != nil {
			logger.Log.Error("Failed to format report", "error", err)
			return
		}

		// Save report
		outputFile := filepath.Join(outputDir, fmt.Sprintf("report.%s", outputFormat))
		err = report.SaveToFile(outputData, outputFile)
		if err != nil {
			logger.Log.Error("Failed to save report", "error", err)
		} else {
			fmt.Printf("Report saved to %s\n", outputFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().String("input", "scan-results.json", "Input file with scan results")
	reportCmd.Flags().String("output-dir", "./reports", "Directory to save the generated report")
}
