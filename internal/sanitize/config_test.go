package sanitize_test

import (
	"os"
	"testing"

	"github.com/envlens/internal/sanitize"
)

const sampleConfigJSON = `{
  "rules": [
    {"exactKeys": ["DB_PASS", "API_KEY"], "placeholder": "[HIDDEN]"},
    {"keyPrefix": "INTERNAL_"}
  ]
}`

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := sanitize.ParseConfig([]byte(sampleConfigJSON))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(cfg.Rules))
	}
}

func TestParseConfig_Invalid(t *testing.T) {
	_, err := sanitize.ParseConfig([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	f, err := os.CreateTemp("", "sanitize-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(sampleConfigJSON)
	f.Close()

	cfg, err := sanitize.LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(cfg.Rules))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := sanitize.LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestConfig_ToSanitizer(t *testing.T) {
	cfg, _ := sanitize.ParseConfig([]byte(sampleConfigJSON))
	s := cfg.ToSanitizer()
	if s == nil {
		t.Fatal("expected non-nil sanitizer")
	}
	out := s.Apply(entries("DB_PASS", "secret", "APP_ENV", "prod"))
	for _, e := range out {
		if e.Key == "DB_PASS" && e.Value != "[HIDDEN]" {
			t.Errorf("expected [HIDDEN], got %q", e.Value)
		}
	}
}
