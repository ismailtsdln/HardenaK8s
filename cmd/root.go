package cmd

import (
	"fmt"
	"os"

	"github.com/ismailtsdln/HardenaK8s/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hardena",
	Short: "HardenaK8s is a Kubernetes security auditing and hardening tool",
	Long: ui.StyleTitle.Render(" HARDENAK8S ") + `

HardenaK8s is a powerful CLI tool designed to audit Kubernetes clusters 
for security misconfigurations and provide actionable hardening recommendations.
It supports CIS Benchmarks, custom policies, and various output formats like JSON/YAML/HTML.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hardena.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "text", "Output format (text, json, yaml, html)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".hardena" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hardena")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
