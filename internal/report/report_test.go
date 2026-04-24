package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/report"
	"github.com/user/envlens/internal/schema"
	"github.com/user/envlens/internal/secrets"
)

func makeSummary() report.Summary {
	return report.Summary{
		Diffs: []differ.DiffEntry{
			{Key: "DB_HOST", Type: differ.Added, NewValue: "localhost"},
			{Key: "API_KEY", Type: differ.Changed, OldValue: "old", NewValue: "new"},
			{Key: "LEGACY", Type: differ.Removed, OldValue: "yes"},
		},
		Findings: []secrets.Finding{
			{Key: "API_KEY", Reason: "high entropy", Severity: "HIGH"},
		},
		Violations: []schema.Violation{
			{Key: "PORT", Message: "missing required key"},
		},
	}
}

func TestWriteText_ContainsSections(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	if err := w.Write(makeSummary()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"ENVLENS REPORT", "[DIFF]", "[SECRETS]", "[SCHEMA]", "DB_HOST", "API_KEY", "PORT"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestWriteText_DiffSymbols(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	_ = w.Write(makeSummary())
	out := buf.String()
	if !strings.Contains(out, "+ DB_HOST") {
		t.Error("expected '+' prefix for added key")
	}
	if !strings.Contains(out, "- LEGACY") {
		t.Error("expected '-' prefix for removed key")
	}
	if !strings.Contains(out, "~ API_KEY") {
		t.Error("expected '~' prefix for changed key")
	}
}

func TestWriteJSON_OutputFormat(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatJSON)
	if err := w.Write(makeSummary()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, "\"diffs\":3") {
		t.Errorf("expected diffs count 3 in JSON, got: %s", out)
	}
}

func TestWriteText_EmptySummary(t *testing.T) {
	var buf bytes.Buffer
	w := report.NewWriter(&buf, report.FormatText)
	if err := w.Write(report.Summary{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 change(s)") {
		t.Error("expected 0 changes in empty summary")
	}
}
