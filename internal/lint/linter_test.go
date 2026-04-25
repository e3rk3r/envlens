package lint_test

import (
	"testing"

	"github.com/user/envlens/internal/lint"
	"github.com/user/envlens/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestLint_Clean(t *testing.T) {
	l := lint.New()
	findings := l.Lint(entries("DB_HOST", "localhost", "DB_PORT", "5432"))
	if len(findings) != 0 {
		t.Fatalf("expected no findings, got %v", findings)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	l := lint.New()
	findings := l.Lint(entries("db_host", "localhost"))
	if len(findings) == 0 {
		t.Fatal("expected warning for lowercase key")
	}
	if findings[0].Severity != lint.SeverityWarning {
		t.Errorf("expected warning severity, got %s", findings[0].Severity)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	l := lint.New()
	findings := l.Lint(entries("DB_HOST", "a", "DB_HOST", "b"))
	found := false
	for _, f := range findings {
		if f.Severity == lint.SeverityError && f.Key == "DB_HOST" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected error finding for duplicate key")
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	l := lint.New()
	findings := l.Lint([]parser.Entry{{Key: "MY_VAR", Value: " value "}})
	if len(findings) == 0 {
		t.Fatal("expected warning for whitespace in value")
	}
	if findings[0].Key != "MY_VAR" {
		t.Errorf("unexpected key %s", findings[0].Key)
	}
}

func TestLint_CustomRule(t *testing.T) {
	noEmptyValue := func(line int, e parser.Entry) []lint.Finding {
		if e.Value == "" {
			return []lint.Finding{{Line: line, Key: e.Key, Message: "empty value", Severity: lint.SeverityInfo}}
		}
		return nil
	}
	l := lint.WithRules(noEmptyValue)
	findings := l.Lint(entries("EMPTY_KEY", "", "GOOD_KEY", "ok"))
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Key != "EMPTY_KEY" {
		t.Errorf("unexpected key %s", findings[0].Key)
	}
}
