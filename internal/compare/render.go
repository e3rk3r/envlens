package compare

import (
	"fmt"
	"io"
	"strings"
)

const (
	absent  = "<absent>"
	maxCell = 20
)

// RenderTable writes a human-readable comparison table to w.
func RenderTable(w io.Writer, m *Matrix) {
	if len(m.Keys) == 0 {
		fmt.Fprintln(w, "(no keys found)")
		return
	}

	// header
	fmt.Fprintf(w, "%-30s", "KEY")
	for _, env := range m.Envs {
		fmt.Fprintf(w, "  %-20s", truncate(env, maxCell))
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 30+len(m.Envs)*22))

	for _, key := range m.Keys {
		marker := " "
		if _, hasMissing := m.Missing[key]; hasMissing {
			marker = "!"
		}
		fmt.Fprintf(w, "%s %-29s", marker, truncate(key, 29))
		for _, env := range m.Envs {
			val := m.Cells[key][env]
			if val == "" {
				val = absent
			}
			fmt.Fprintf(w, "  %-20s", truncate(val, maxCell))
		}
		fmt.Fprintln(w)
	}
}

// RenderMissingReport writes a focused report of keys absent in any environment.
func RenderMissingReport(w io.Writer, m *Matrix) {
	if len(m.Missing) == 0 {
		fmt.Fprintln(w, "No missing keys across environments.")
		return
	}
	fmt.Fprintln(w, "Missing keys:")
	for _, key := range m.Keys {
		envs, ok := m.Missing[key]
		if !ok {
			continue
		}
		fmt.Fprintf(w, "  %-30s absent in: %s\n", key, strings.Join(envs, ", "))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
