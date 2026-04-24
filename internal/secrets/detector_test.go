package secrets

import (
	"testing"
)

func TestScan_NoFindings(t *testing.T) {
	d := NewDetector()
	entries := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
		"DEBUG":    "true",
	}
	findings := d.Scan(entries)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings, got %d", len(findings))
	}
}

func TestScan_DetectsHighEntropySecret(t *testing.T) {
	d := NewDetector()
	entries := map[string]string{
		"API_SECRET": "aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890",
	}
	findings := d.Scan(entries)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "high-entropy-secret" {
		t.Errorf("expected rule 'high-entropy-secret', got %q", findings[0].RuleID)
	}
}

func TestScan_DetectsAWSKey(t *testing.T) {
	d := NewDetector()
	entries := map[string]string{
		"CLOUD_KEY": "AKIAIOSFODNN7EXAMPLE",
	}
	findings := d.Scan(entries)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "aws-access-key" {
		t.Errorf("expected rule 'aws-access-key', got %q", findings[0].RuleID)
	}
}

func TestScan_DetectsJWT(t *testing.T) {
	d := NewDetector()
	entries := map[string]string{
		"AUTH_TOKEN": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0",
	}
	findings := d.Scan(entries)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].RuleID != "generic-bearer-token" {
		t.Errorf("expected rule 'generic-bearer-token', got %q", findings[0].RuleID)
	}
}

func TestScan_SkipsEmptyValues(t *testing.T) {
	d := NewDetector()
	entries := map[string]string{
		"API_SECRET": "",
		"PASSWORD":   "   ",
	}
	findings := d.Scan(entries)
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for empty values, got %d", len(findings))
	}
}
