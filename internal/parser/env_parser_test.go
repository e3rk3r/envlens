package parser

import (
	"strings"
	"testing"
)

func TestParse_BasicKeyValue(t *testing.T) {
	input := "APP_ENV=production\nPORT=8080\n"
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(ef.Entries))
	}
	if ef.Raw["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", ef.Raw["APP_ENV"])
	}
	if ef.Raw["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", ef.Raw["PORT"])
	}
}

func TestParse_SkipsCommentsAndBlanks(t *testing.T) {
	input := "# this is a comment\n\nFOO=bar\n"
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "FOO" {
		t.Errorf("expected key FOO, got %q", ef.Entries[0].Key)
	}
}

func TestParse_QuotedValues(t *testing.T) {
	input := `DB_URL="postgres://localhost/mydb"` + "\n" + `SECRET='mysecret'` + "\n"
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Raw["DB_URL"] != "postgres://localhost/mydb" {
		t.Errorf("unexpected DB_URL: %q", ef.Raw["DB_URL"])
	}
	if ef.Raw["SECRET"] != "mysecret" {
		t.Errorf("unexpected SECRET: %q", ef.Raw["SECRET"])
	}
}

func TestParse_InlineComment(t *testing.T) {
	input := "HOST=localhost #local only\n"
	ef, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Raw["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", ef.Raw["HOST"])
	}
	if ef.Entries[0].Comment != "local only" {
		t.Errorf("expected comment 'local only', got %q", ef.Entries[0].Comment)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	input := "INVALID_LINE_WITHOUT_EQUALS\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_EmptyKey(t *testing.T) {
	input := "=value\n"
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}
