package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ismailtsdln/HardenaK8s/internal/k8s"
	"github.com/ismailtsdln/HardenaK8s/internal/logger"
	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"github.com/ismailtsdln/HardenaK8s/internal/report"
	"github.com/ismailtsdln/HardenaK8s/internal/ui"
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

		fmt.Println(ui.StyleHeader.Render("Starting Security Scan..."))

		client, err := k8s.NewClient()
		if err != nil {
			fmt.Println(ui.Error("Failed to initialize Kubernetes client: " + err.Error()))
			os.Exit(1)
		}

		ctx := context.Background()
		fmt.Println(ui.Info("Checking cluster connectivity..."))
		if err := client.CheckConnectivity(ctx); err != nil {
			fmt.Println(ui.Error("Could not connect to Kubernetes cluster: " + err.Error()))
			os.Exit(1)
		}
		fmt.Println(ui.Success("Connected to cluster."))

		fmt.Println(ui.Info("Auditing resources..."))
		engine := policy.NewEngine(client)
		result, err := engine.Run(ctx, namespace)
		if err != nil {
			fmt.Println(ui.Error("Scan failed: " + err.Error()))
			os.Exit(1)
		}

		logger.Log.Info("Scan completed", "issues_found", result.Stats.TotalIssues)

		outputFormat := viper.GetString("output")

		if outputFormat == "text" {
			renderTable(result)
		} else {
			formatter, err := report.GetFormatter(outputFormat)
			if err != nil {
				fmt.Println(ui.Warning("Invalid output format, defaulting to JSON"))
				formatter = &report.JSONFormatter{}
				outputFormat = "json"
			}

			data, err := formatter.Format(result)
			if err != nil {
				fmt.Println(ui.Error("Failed to format report: " + err.Error()))
				return
			}

			outputFile := fmt.Sprintf("scan-results.%s", outputFormat)
			err = report.SaveToFile(data, outputFile)
			if err != nil {
				fmt.Println(ui.Error("Failed to save report: " + err.Error()))
			} else {
				fmt.Println(ui.Success("Report saved to " + outputFile))
			}
		}
	},
}

func renderTable(result *policy.Result) {
	if len(result.Issues) == 0 {
		fmt.Println("\n" + ui.Success("No security issues found! Your cluster is hardened. 🛡️"))
		return
	}

	fmt.Println(ui.StyleHeader.Render("\nSecurity Findings Summary"))

	for _, issue := range result.Issues {
		var sevStyle = ui.StyleInfo
		switch issue.Severity {
		case policy.SeverityCritical:
			sevStyle = ui.StyleCritical
		case policy.SeverityHigh:
			sevStyle = ui.StyleError
		case policy.SeverityMedium:
			sevStyle = ui.StyleWarn
		}

		fmt.Printf("[%s] %s\n", sevStyle.Render(string(issue.Severity)), ui.StyleHeader.Render(issue.Title))
		fmt.Printf("   Resource: %s/%s\n", issue.Namespace, issue.Resource)
		fmt.Printf("   Details:  %s\n", issue.Description)
		fmt.Printf("   Fix:      %s\n\n", ui.StyleSuccess.Render(issue.Remediation))
	}

	fmt.Println(ui.StyleHeader.Render("Scan Statistics"))
	fmt.Printf("Total Issues:    %d\n", result.Stats.TotalIssues)
	for sev, count := range result.Stats.SeverityCount {
		fmt.Printf("%-15s %d\n", sev+":", count)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().String("namespace", "", "Scan a specific namespace")
	scanCmd.Flags().Bool("all-namespaces", true, "Scan all namespaces")
}
