package schema

import (
	"testing"
)

func TestValidate_AllPresentAndValid(t *testing.T) {
	schema := Schema{
		"APP_ENV": {Required: true, Allowed: []string{"development", "production", "staging"}},
		"PORT":    {Required: true, Pattern: `^\d+$`},
	}
	env := map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	}
	violations := Validate(schema, env)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	schema := Schema{
		"DATABASE_URL": {Required: true},
	}
	env := map[string]string{}
	violations := Validate(schema, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected violation for DATABASE_URL, got %s", violations[0].Key)
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	schema := Schema{
		"PORT": {Required: true, Pattern: `^\d+$`},
	}
	env := map[string]string{"PORT": "not-a-number"}
	violations := Validate(schema, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_AllowedValueViolation(t *testing.T) {
	schema := Schema{
		"LOG_LEVEL": {Allowed: []string{"debug", "info", "warn", "error"}},
	}
	env := map[string]string{"LOG_LEVEL": "verbose"}
	violations := Validate(schema, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_OptionalKeyAbsent(t *testing.T) {
	schema := Schema{
		"OPTIONAL_FEATURE": {Required: false, Allowed: []string{"true", "false"}},
	}
	env := map[string]string{}
	violations := Validate(schema, env)
	if len(violations) != 0 {
		t.Errorf("expected no violations for absent optional key, got %v", violations)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	schema := Schema{
		"APP_ENV":      {Required: true},
		"PORT":         {Required: true, Pattern: `^\d+$`},
		"DATABASE_URL": {Required: true},
	}
	env := map[string]string{
		"PORT": "abc",
	}
	violations := Validate(schema, env)
	if len(violations) != 3 {
		t.Errorf("expected 3 violations, got %d: %v", len(violations), violations)
	}
}
