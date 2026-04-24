// Package main is the entry point for the envlens CLI tool.
// envlens helps developers diff, audit, and sanitize .env files
// across environments, with secret detection and schema validation.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envlens",
	Short: "Audit, diff, and sanitize .env files",
	Long: `envlens is a tool to diff, audit, and sanitize .env files across
environments. It supports secret detection, schema validation, and
produces reports in text or JSON format.`,
}

// diffCmd compares two .env files and reports differences.
var diffCmd = &cobra.Command{
	Use:   "diff <base> <target>",
	Short: "Diff two .env files",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

// auditCmd scans a .env file for secrets and schema violations.
var auditCmd = &cobra.Command{
	Use:   "audit <file>",
	Short: "Audit a .env file for secrets and schema violations",
	Args:  cobra.ExactArgs(1),
	RunE:  runAudit,
}

// sanitizeCmd applies sanitization rules to a .env file.
var sanitizeCmd = &cobra.Command{
	Use:   "sanitize <file>",
	Short: "Sanitize a .env file using configured rules",
	Args:  cobra.ExactArgs(1),
	RunE:  runSanitize,
}

func init() {
	// diff flags
	diffCmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	diffCmd.Flags().StringP("schema", "s", "", "Path to schema YAML file for validation")

	// audit flags
	auditCmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	auditCmd.Flags().StringP("schema", "s", "", "Path to schema YAML file for validation")
	auditCmd.Flags().BoolP("redact", "r", false, "Redact secret values in output")

	// sanitize flags
	sanitizeCmd.Flags().StringP("config", "c", "", "Path to sanitizer config YAML file")
	sanitizeCmd.Flags().StringP("output", "o", "-", "Output file path (default: stdout)")

	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(auditCmd)
	rootCmd.AddCommand(sanitizeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
