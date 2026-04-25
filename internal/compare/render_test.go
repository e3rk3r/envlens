package compare_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/compare"
	"github.com/user/envlens/internal/parser"
)

func TestRenderTable_ContainsHeaders(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("PORT", "3000"),
		"prod": entries("PORT", "8080"),
	}
	m := compare.Compare(envs)
	var buf bytes.Buffer
	compare.RenderTable(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "KEY") {
		t.Error("expected header KEY in output")
	}
	if !strings.Contains(out, "dev") {
		t.Error("expected env name 'dev' in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected key PORT in output")
	}
}

func TestRenderTable_ShowsAbsentMarker(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("ONLY_DEV", "yes"),
		"prod": entries("OTHER", "no"),
	}
	m := compare.Compare(envs)
	var buf bytes.Buffer
	compare.RenderTable(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "<absent>") {
		t.Error("expected <absent> marker for missing key")
	}
	if !strings.Contains(out, "!") {
		t.Error("expected '!' marker for row with missing key")
	}
}

func TestRenderTable_EmptyMatrix(t *testing.T) {
	m := compare.Compare(map[string][]parser.Entry{})
	var buf bytes.Buffer
	compare.RenderTable(&buf, m)
	if !strings.Contains(buf.String(), "no keys") {
		t.Error("expected 'no keys' message for empty matrix")
	}
}

func TestRenderMissingReport_NoMissing(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("A", "1"),
		"prod": entries("A", "2"),
	}
	m := compare.Compare(envs)
	var buf bytes.Buffer
	compare.RenderMissingReport(&buf, m)
	if !strings.Contains(buf.String(), "No missing") {
		t.Error("expected no-missing message")
	}
}

func TestRenderMissingReport_HasMissing(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("A", "1", "B", "2"),
		"prod": entries("A", "1"),
	}
	m := compare.Compare(envs)
	var buf bytes.Buffer
	compare.RenderMissingReport(&buf, m)
	out := buf.String()
	if !strings.Contains(out, "B") {
		t.Error("expected key B in missing report")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected env 'prod' in missing report")
	}
}
