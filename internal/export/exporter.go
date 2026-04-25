package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Format represents the output format for exported env entries.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Exporter writes env entries to an output stream in a given format.
type Exporter struct {
	format Format
}

// New creates a new Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatJSON, FormatShell:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write exports the given entries to w in the configured format.
func (e *Exporter) Write(entries []parser.Entry, w io.Writer) error {
	switch e.format {
	case FormatDotenv:
		return writeDotenv(entries, w)
	case FormatJSON:
		return writeJSON(entries, w)
	case FormatShell:
		return writeShell(entries, w)
	}
	return nil
}

func writeDotenv(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		line := fmt.Sprintf("%s=%s\n", e.Key, quoteIfNeeded(e.Value))
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(entries []parser.Entry, w io.Writer) error {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ordered := make([]map[string]string, 0, len(keys))
	for _, k := range keys {
		ordered = append(ordered, map[string]string{"key": k, "value": m[k]})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(ordered)
}

func writeShell(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		line := fmt.Sprintf("export %s=%s\n", e.Key, quoteIfNeeded(e.Value))
		if _, err := io.WriteString(w, line); err != nil {
			return err
		}
	}
	return nil
}

func quoteIfNeeded(v string) string {
	if strings.ContainsAny(v, " \t\n#") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
