package lint

import (
	"testing"

	"github.com/user/envlens/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestRuleUpperCaseKeys_Clean(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost", "API_KEY", "secret")
	findings := RuleUpperCaseKeys(entries)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestRuleUpperCaseKeys_Violation(t *testing.T) {
	entries := makeEntries("db_host", "localhost", "Api_Key", "secret")
	findings := RuleUpperCaseKeys(entries)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	for _, f := range findings {
		if f.Rule != "upper-case-keys" {
			t.Errorf("unexpected rule: %s", f.Rule)
		}
	}
}

func TestRuleNoDuplicateKeys_Clean(t *testing.T) {
	entries := makeEntries("FOO", "1", "BAR", "2")
	findings := RuleNoDuplicateKeys(entries)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestRuleNoDuplicateKeys_Duplicate(t *testing.T) {
	entries := []parser.Entry{
		{Key: "FOO", Value: "1"},
		{Key: "FOO", Value: "2"},
		{Key: "BAR", Value: "3"},
	}
	findings := RuleNoDuplicateKeys(entries)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %s", findings[0].Key)
	}
}

func TestRuleNoWhitespaceValues_Clean(t *testing.T) {
	entries := makeEntries("HOST", "localhost", "PORT", "5432")
	findings := RuleNoWhitespaceValues(entries)
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(findings))
	}
}

func TestRuleNoWhitespaceValues_LeadingSpace(t *testing.T) {
	entries := makeEntries("HOST", " localhost")
	findings := RuleNoWhitespaceValues(entries)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "no-whitespace-values" {
		t.Errorf("unexpected rule: %s", findings[0].Rule)
	}
}
