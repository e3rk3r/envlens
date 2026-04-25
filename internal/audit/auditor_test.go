package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlens/internal/audit"
	"github.com/user/envlens/internal/secrets"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return filepath.Clean(f.Name())
}

func TestRun_NoDiffsNoIssues(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")
	target := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")

	detector := secrets.NewDetector(nil)
	a := audit.NewAuditor(detector, nil)

	result, err := a.Run(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(result.Diffs))
	}
	if result.HasIssues() {
		t.Error("expected no issues")
	}
}

func TestRun_DetectsDiff(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=staging\n")
	target := writeTempEnv(t, "APP_ENV=production\nNEW_KEY=value\n")

	detector := secrets.NewDetector(nil)
	a := audit.NewAuditor(detector, nil)

	result, err := a.Run(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diffs) == 0 {
		t.Error("expected diffs but got none")
	}
}

func TestRun_DetectsSecret(t *testing.T) {
	base := writeTempEnv(t, "APP_ENV=production\n")
	target := writeTempEnv(t, "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\n")

	detector := secrets.NewDetector(nil)
	a := audit.NewAuditor(detector, nil)

	result, err := a.Run(base, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasIssues() {
		t.Error("expected secret finding but got none")
	}
}

func TestRun_InvalidBaseFile(t *testing.T) {
	detector := secrets.NewDetector(nil)
	a := audit.NewAuditor(detector, nil)

	_, err := a.Run("/nonexistent/base.env", "/nonexistent/target.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
