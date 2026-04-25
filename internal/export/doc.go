// Package export provides functionality to write parsed .env entries
// to various output formats including dotenv, JSON, and shell (export statements).
//
// Supported formats:
//
//   - dotenv  — standard KEY=VALUE format, values quoted when necessary
//   - json    — JSON array of {"key", "value"} objects, sorted alphabetically
//   - shell   — shell-compatible `export KEY=VALUE` lines
//
// Example usage:
//
//	exporter, err := export.New(export.FormatJSON)
//	if err != nil { ... }
//	exporter.Write(entries, os.Stdout)
package export
