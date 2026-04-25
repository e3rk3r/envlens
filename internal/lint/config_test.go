package lint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/lint"
	"github.com/user/envlens/internal/parser"
)

func TestParseConfig_Valid(t *testing.T) {
	data := []byte(`{"require_upper_case":true,"forbid_duplicate_keys":false,"forbid_whitespace":true}`)
	cfg, err := lint.ParseConfig(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.RequireUpperCase {
		t.Error("expected RequireUpperCase=true")
	}
	if cfg.ForbidDuplicateKeys {
		t.Error("expected ForbidDuplicateKeys=false")
	}
}

func TestParseConfig_Invalid(t *testing.T) {
	_, err := lint.ParseConfig([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lint.json")
	content := `{"require_upper_case":false,"forbid_duplicate_keys":true,"forbid_whitespace":false}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := lint.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.RequireUpperCase {
		t.Error("expected RequireUpperCase=false")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := lint.LoadConfig("/nonexistent/lint.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestConfig_ToLinter_OnlyUpperCase(t *testing.T) {
	cfg := lint.Config{
		RequireUpperCase:    true,
		ForbidDuplicateKeys: false,
		ForbidWhitespace:    false,
	}
	l := cfg.ToLinter()
	// lowercase key should trigger a finding
	findings := l.Lint([]parser.Entry{{Key: "lower", Value: "v"}})
	if len(findings) == 0 {
		t.Fatal("expected at least one finding for lowercase key")
	}
	// whitespace should NOT trigger (rule disabled)
	findings2 := l.Lint([]parser.Entry{{Key: "GOOD", Value: " spaced "}})
	if len(findings2) != 0 {
		t.Fatalf("expected no findings when whitespace rule disabled, got %v", findings2)
	}
}
