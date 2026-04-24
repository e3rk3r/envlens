package sanitize_test

import (
	"testing"

	"github.com/envlens/internal/parser"
	"github.com/envlens/internal/sanitize"
)

func entries(kvs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestApply_NoRules(t *testing.T) {
	s := sanitize.NewSanitizer(nil)
	in := entries("APP_ENV", "production", "DB_PASS", "secret")
	out := s.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestApply_DropByExactKey(t *testing.T) {
	s := sanitize.NewSanitizer([]sanitize.Rule{
		{ExactKeys: []string{"DB_PASS"}},
	})
	out := s.Apply(entries("APP_ENV", "prod", "DB_PASS", "secret"))
	if len(out) != 1 || out[0].Key != "APP_ENV" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestApply_DropByKeyPrefix(t *testing.T) {
	s := sanitize.NewSanitizer([]sanitize.Rule{
		{KeyPrefix: "INTERNAL_"},
	})
	out := s.Apply(entries("INTERNAL_TOKEN", "abc", "PUBLIC_URL", "https://x.com"))
	if len(out) != 1 || out[0].Key != "PUBLIC_URL" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestApply_DropByKeySuffix(t *testing.T) {
	s := sanitize.NewSanitizer([]sanitize.Rule{
		{KeySuffix: "_SECRET"},
	})
	out := s.Apply(entries("API_SECRET", "xyz", "APP_NAME", "myapp"))
	if len(out) != 1 || out[0].Key != "APP_NAME" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

func TestApply_RedactKey(t *testing.T) {
	s := sanitize.NewSanitizer([]sanitize.Rule{
		{RedactKeys: []string{"DB_PASS"}, Placeholder: "[REDACTED]"},
	})
	out := s.Apply(entries("DB_PASS", "supersecret", "APP_ENV", "prod"))
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "DB_PASS" && e.Value != "[REDACTED]" {
			t.Errorf("expected redacted value, got %q", e.Value)
		}
	}
}

func TestApply_DefaultPlaceholder(t *testing.T) {
	s := sanitize.NewSanitizer([]sanitize.Rule{
		{RedactKeys: []string{"SECRET_KEY"}},
	})
	out := s.Apply(entries("SECRET_KEY", "abc123"))
	if out[0].Value != "***" {
		t.Errorf("expected *** placeholder, got %q", out[0].Value)
	}
}
