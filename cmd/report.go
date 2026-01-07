package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"github.com/ismailtsdln/HardenaK8s/internal/report"
	"github.com/ismailtsdln/HardenaK8s/internal/ui"
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

		fmt.Println(ui.StyleHeader.Render("Generating Security Report..."))
		fmt.Println(ui.Info(fmt.Sprintf("Input:  %s", inputFile)))
		fmt.Println(ui.Info(fmt.Sprintf("Format: %s", outputFormat)))

		// Create output directory if it doesn't exist
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Println(ui.Error("Failed to create output directory: " + err.Error()))
			return
		}

		// Read input results
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Println(ui.Error("Failed to read input file: " + err.Error()))
			return
		}

		var result policy.Result
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Println(ui.Error("Failed to parse scan results: " + err.Error()))
			return
		}

		// Safety check for unmarshaled data
		if result.Stats.SeverityCount == nil {
			result.Stats.SeverityCount = make(map[policy.Severity]int)
		}

		// Format report
		formatter, err := report.GetFormatter(outputFormat)
		if err != nil {
			fmt.Println(ui.Warning("Invalid output format, defaulting to JSON"))
			formatter = &report.JSONFormatter{}
			outputFormat = "json"
		}

		outputData, err := formatter.Format(&result)
		if err != nil {
			fmt.Println(ui.Error("Failed to format report: " + err.Error()))
			return
		}

		// Save report
		outputFile := filepath.Join(outputDir, fmt.Sprintf("report.%s", outputFormat))
		err = report.SaveToFile(outputData, outputFile)
		if err != nil {
			fmt.Println(ui.Error("Failed to save report: " + err.Error()))
		} else {
			fmt.Println(ui.Success("Report successfully saved to " + outputFile))
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().String("input", "scan-results.json", "Input file with scan results")
	reportCmd.Flags().String("output-dir", "./reports", "Directory to save the generated report")
}
