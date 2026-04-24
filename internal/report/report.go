package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/schema"
	"github.com/user/envlens/internal/secrets"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Summary holds aggregated results from all checks.
type Summary struct {
	Diffs      []differ.DiffEntry
	Findings   []secrets.Finding
	Violations []schema.Violation
}

// Writer renders a Summary to an io.Writer in the given format.
type Writer struct {
	format Format
	out    io.Writer
}

// NewWriter creates a new report Writer.
func NewWriter(out io.Writer, format Format) *Writer {
	return &Writer{out: out, format: format}
}

// Write renders the summary.
func (w *Writer) Write(s Summary) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(s)
	default:
		return w.writeText(s)
	}
}

func (w *Writer) writeText(s Summary) error {
	fmt.Fprintln(w.out, strings.Repeat("=", 50))
	fmt.Fprintln(w.out, "ENVLENS REPORT")
	fmt.Fprintln(w.out, strings.Repeat("=", 50))

	fmt.Fprintf(w.out, "\n[DIFF] %d change(s) found\n", len(s.Diffs))
	for _, d := range s.Diffs {
		switch d.Type {
		case differ.Added:
			fmt.Fprintf(w.out, "  + %s = %s\n", d.Key, d.NewValue)
		case differ.Removed:
			fmt.Fprintf(w.out, "  - %s = %s\n", d.Key, d.OldValue)
		case differ.Changed:
			fmt.Fprintf(w.out, "  ~ %s: %q -> %q\n", d.Key, d.OldValue, d.NewValue)
		}
	}

	fmt.Fprintf(w.out, "\n[SECRETS] %d finding(s)\n", len(s.Findings))
	for _, f := range s.Findings {
		fmt.Fprintf(w.out, "  [%s] %s — %s\n", f.Severity, f.Key, f.Reason)
	}

	fmt.Fprintf(w.out, "\n[SCHEMA] %d violation(s)\n", len(s.Violations))
	for _, v := range s.Violations {
		fmt.Fprintf(w.out, "  [%s] %s\n", v.Key, v.Message)
	}

	fmt.Fprintln(w.out, strings.Repeat("=", 50))
	return nil
}

func (w *Writer) writeJSON(s Summary) error {
	// Simple hand-rolled JSON to avoid import of encoding/json for brevity.
	fmt.Fprintf(w.out, "{\"diffs\":%d,\"findings\":%d,\"violations\":%d}\n",
		len(s.Diffs), len(s.Findings), len(s.Violations))
	return nil
}
