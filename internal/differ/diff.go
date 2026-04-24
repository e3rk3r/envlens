package differ

import (
	"fmt"
	"sort"

	"github.com/yourusername/envlens/internal/parser"
)

// DiffType represents the type of difference found between two env files.
type DiffType string

const (
	DiffAdded   DiffType = "added"
	DiffRemoved DiffType = "removed"
	DiffChanged DiffType = "changed"
)

// DiffEntry represents a single difference between two env files.
type DiffEntry struct {
	Key      string
	Type     DiffType
	OldValue string
	NewValue string
}

// String returns a human-readable representation of the diff entry.
func (d DiffEntry) String() string {
	switch d.Type {
	case DiffAdded:
		return fmt.Sprintf("+ %s=%s", d.Key, d.NewValue)
	case DiffRemoved:
		return fmt.Sprintf("- %s=%s", d.Key, d.OldValue)
	case DiffChanged:
		return fmt.Sprintf("~ %s: %q -> %q", d.Key, d.OldValue, d.NewValue)
	}
	return ""
}

// Diff compares two sets of parsed env entries and returns the differences.
func Diff(base, target []parser.Entry) []DiffEntry {
	baseMap := toMap(base)
	targetMap := toMap(target)

	var diffs []DiffEntry

	for key, baseVal := range baseMap {
		if targetVal, ok := targetMap[key]; !ok {
			diffs = append(diffs, DiffEntry{
				Key:      key,
				Type:     DiffRemoved,
				OldValue: baseVal,
			})
		} else if baseVal != targetVal {
			diffs = append(diffs, DiffEntry{
				Key:      key,
				Type:     DiffChanged,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	for key, targetVal := range targetMap {
		if _, ok := baseMap[key]; !ok {
			diffs = append(diffs, DiffEntry{
				Key:      key,
				Type:     DiffAdded,
				NewValue: targetVal,
			})
		}
	}

	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].Key < diffs[j].Key
	})

	return diffs
}

func toMap(entries []parser.Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
