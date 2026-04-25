package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envlens/internal/export"
	"github.com/user/envlens/internal/parser"
)

var exportFormat string

var exportCmd = &cobra.Command{
	Use:   "export <file>",
	Short: "Export a .env file in a different format (dotenv, json, shell)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}
		defer f.Close()

		entries, err := parser.Parse(f)
		if err != nil {
			return fmt.Errorf("parsing file: %w", err)
		}

		exporter, err := export.New(export.Format(exportFormat))
		if err != nil {
			return err
		}

		return exporter.Write(entries, os.Stdout)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv",
		"Output format: dotenv, json, or shell")
	rootCmd.AddCommand(exportCmd)
}
