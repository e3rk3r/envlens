package differ_test

import (
	"testing"

	"github.com/yourusername/envlens/internal/differ"
	"github.com/yourusername/envlens/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var result []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		result = append(result, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return result
}

func TestDiff_NoDifferences(t *testing.T) {
	base := entries("FOO", "bar", "BAZ", "qux")
	target := entries("FOO", "bar", "BAZ", "qux")
	diffs := differ.Diff(base, target)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiff_AddedKey(t *testing.T) {
	base := entries("FOO", "bar")
	target := entries("FOO", "bar", "NEW_KEY", "new_val")
	diffs := differ.Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != differ.DiffAdded || diffs[0].Key != "NEW_KEY" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := entries("FOO", "bar", "OLD_KEY", "old_val")
	target := entries("FOO", "bar")
	diffs := differ.Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != differ.DiffRemoved || diffs[0].Key != "OLD_KEY" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
}

func TestDiff_ChangedValue(t *testing.T) {
	base := entries("DB_HOST", "localhost")
	target := entries("DB_HOST", "prod.db.example.com")
	diffs := differ.Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if d.Type != differ.DiffChanged || d.OldValue != "localhost" || d.NewValue != "prod.db.example.com" {
		t.Errorf("unexpected diff: %+v", d)
	}
}

func TestDiff_SortedOutput(t *testing.T) {
	base := entries("Z_KEY", "1", "A_KEY", "2")
	target := entries("Z_KEY", "changed", "A_KEY", "changed")
	diffs := differ.Diff(base, target)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	if diffs[0].Key != "A_KEY" || diffs[1].Key != "Z_KEY" {
		t.Errorf("diffs not sorted: %v, %v", diffs[0].Key, diffs[1].Key)
	}
}

func TestDiffEntry_String(t *testing.T) {
	cases := []struct {
		entry    differ.DiffEntry
		expected string
	}{
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffAdded, NewValue: "bar"}, "+ FOO=bar"},
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffRemoved, OldValue: "bar"}, "- FOO=bar"},
		{differ.DiffEntry{Key: "FOO", Type: differ.DiffChanged, OldValue: "old", NewValue: "new"}, `~ FOO: "old" -> "new"`},
	}
	for _, tc := range cases {
		got := tc.entry.String()
		if got != tc.expected {
			t.Errorf("String() = %q, want %q", got, tc.expected)
		}
	}
}
