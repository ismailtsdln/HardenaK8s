package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ismailtsdln/HardenaK8s/internal/k8s"
	"github.com/ismailtsdln/HardenaK8s/internal/logger"
	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"github.com/ismailtsdln/HardenaK8s/internal/report"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan Kubernetes cluster for security issues",
	Long: `The scan command audits the Kubernetes cluster against a set of 
predefined and custom security policies. It checks for common misconfigurations,
RBAC issues, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

		if allNamespaces {
			namespace = ""
		}

		logger.Log.Info("Starting scan", "namespace", namespace)

		client, err := k8s.NewClient()
		if err != nil {
			logger.Log.Error("Failed to initialize Kubernetes client", "error", err)
			os.Exit(1)
		}

		ctx := context.Background()
		engine := policy.NewEngine(client)
		result, err := engine.Run(ctx, namespace)
		if err != nil {
			logger.Log.Error("Scan failed", "error", err)
			os.Exit(1)
		}

		logger.Log.Info("Scan completed", "issues_found", result.Stats.TotalIssues)

		outputFormat := viper.GetString("output")
		formatter, err := report.GetFormatter(outputFormat)
		if err != nil {
			logger.Log.Warn("Invalid output format, defaulting to JSON", "format", outputFormat)
			formatter = &report.JSONFormatter{}
		}

		data, err := formatter.Format(result)
		if err != nil {
			logger.Log.Error("Failed to format report", "error", err)
			return
		}

		// Print to stdout or save to file
		outputFile := fmt.Sprintf("scan-results.%s", outputFormat)
		if outputFormat == "text" {
			// Fallback if text is requested but not implemented yet
			fmt.Println(string(data))
		} else {
			err = report.SaveToFile(data, outputFile)
			if err != nil {
				logger.Log.Error("Failed to save report", "error", err)
			} else {
				fmt.Printf("Report saved to %s\n", outputFile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().String("namespace", "", "Scan a specific namespace")
	scanCmd.Flags().Bool("all-namespaces", true, "Scan all namespaces")
}
