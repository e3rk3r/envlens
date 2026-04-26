package lint

import (
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Finding represents a single lint violation.
type Finding struct {
	Key     string
	Rule    string
	Message string
}

// Rule is a function that inspects a slice of entries and returns findings.
type Rule func(entries []parser.Entry) []Finding

// RuleUpperCaseKeys flags any key that contains lowercase letters.
func RuleUpperCaseKeys(entries []parser.Entry) []Finding {
	var findings []Finding
	for _, e := range entries {
		if e.Key != strings.ToUpper(e.Key) {
			findings = append(findings, Finding{
				Key:     e.Key,
				Rule:    "upper-case-keys",
				Message: "key should be upper-case: " + e.Key,
			})
		}
	}
	return findings
}

// RuleNoDuplicateKeys flags keys that appear more than once.
func RuleNoDuplicateKeys(entries []parser.Entry) []Finding {
	seen := make(map[string]int)
	for _, e := range entries {
		seen[e.Key]++
	}
	var findings []Finding
	for key, count := range seen {
		if count > 1 {
			findings = append(findings, Finding{
				Key:     key,
				Rule:    "no-duplicate-keys",
				Message: "duplicate key detected: " + key,
			})
		}
	}
	return findings
}

// RuleNoWhitespaceValues flags entries whose value contains leading or trailing whitespace.
func RuleNoWhitespaceValues(entries []parser.Entry) []Finding {
	var findings []Finding
	for _, e := range entries {
		if e.Value != strings.TrimSpace(e.Value) {
			findings = append(findings, Finding{
				Key:     e.Key,
				Rule:    "no-whitespace-values",
				Message: "value for key has leading/trailing whitespace: " + e.Key,
			})
		}
	}
	return findings
}
