package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ismailtsdln/HardenaK8s/internal/policy"
	"github.com/ismailtsdln/HardenaK8s/internal/ui"
	"github.com/spf13/cobra"
)

// fixCmd represents the fix command
var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Apply security fixes and hardening",
	Long: `The fix command attempts to automatically remediate identified 
security vulnerabilities or misconfigurations where possible.`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		inputFile, _ := cmd.Flags().GetString("input")

		fmt.Println(ui.StyleHeader.Render("Starting Security Hardening..."))

		if dryRun {
			fmt.Println(ui.Info("Running in Dry Run mode. No changes will be applied."))
		}

		// Read scan results
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Println(ui.Error("Failed to read scan results: " + err.Error()))
			return
		}

		var result policy.Result
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Println(ui.Error("Failed to parse scan results: " + err.Error()))
			return
		}

		if len(result.Issues) == 0 {
			fmt.Println(ui.Success("No issues to fix. Cluster is already hardened."))
			return
		}

		fmt.Println(ui.Info(fmt.Sprintf("Analyzing %d issues...", len(result.Issues))))

		for _, issue := range result.Issues {
			fmt.Printf("\nIssue: %s (%s)\n", ui.StyleHeader.Render(issue.Title), issue.ID)
			fmt.Printf("Action: %s\n", ui.StyleSuccess.Render(issue.Remediation))
		}

		fmt.Println("\n" + ui.Warning("Automated remediation is currently under development."))
		fmt.Println(ui.Info("Please apply the above changes manually to your manifests."))

		if dryRun {
			fmt.Println("\n" + ui.Success("Dry run completed."))
		}
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)

	fixCmd.Flags().Bool("dry-run", true, "Show what would be changed without applying")
	fixCmd.Flags().String("input", "scan-results.json", "Input file with scan results")
}
