package export

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/parser"
)

func entries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "s3cr3t value"},
	}
}

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatDotenv, FormatJSON, FormatShell} {
		e, err := New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
		if e == nil {
			t.Errorf("expected non-nil exporter for format %q", f)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("yaml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_Dotenv(t *testing.T) {
	e, _ := New(FormatDotenv)
	var buf bytes.Buffer
	if err := e.Write(entries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, `DB_PASS=`) {
		t.Errorf("expected DB_PASS in output, got:\n%s", out)
	}
}

func TestWrite_Shell(t *testing.T) {
	e, _ := New(FormatShell)
	var buf bytes.Buffer
	if err := e.Write(entries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export APP_ENV=production") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	e, _ := New(FormatJSON)
	var buf bytes.Buffer
	if err := e.Write(entries(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) || !strings.Contains(out, `"value"`) {
		t.Errorf("expected JSON keys in output, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in JSON output, got:\n%s", out)
	}
}

func TestWrite_DotenvQuotesSpaces(t *testing.T) {
	e, _ := New(FormatDotenv)
	var buf bytes.Buffer
	input := []parser.Entry{{Key: "MSG", Value: "hello world"}}
	if err := e.Write(input, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value for space-containing string, got:\n%s", out)
	}
}
