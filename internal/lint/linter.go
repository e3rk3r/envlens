// Package lint provides .env file linting rules such as naming conventions,
// duplicate key detection, and value format checks.
package lint

import (
	"fmt"
	"strings"

	"github.com/user/envlens/internal/parser"
)

// Severity represents the importance of a lint finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Finding represents a single lint issue found in the env file.
type Finding struct {
	Line     int
	Key      string
	Message  string
	Severity Severity
}

// Linter checks env entries against a set of rules.
type Linter struct {
	rules []Rule
}

// Rule is a function that inspects an entry and returns findings.
type Rule func(line int, entry parser.Entry) []Finding

// New creates a Linter with the default rule set.
func New() *Linter {
	return &Linter{
		rules: []Rule{
			RuleUpperCaseKeys,
			RuleNoDuplicateKeys,
			RuleNoWhitespaceInValue,
		},
	}
}

// WithRules creates a Linter with a custom rule set.
func WithRules(rules ...Rule) *Linter {
	return &Linter{rules: rules}
}

// Lint runs all rules against the provided entries and returns all findings.
func (l *Linter) Lint(entries []parser.Entry) []Finding {
	var findings []Finding
	seen := map[string]int{}

	for i, entry := range entries {
		line := i + 1
		if prev, ok := seen[entry.Key]; ok {
			findings = append(findings, Finding{
				Line:     line,
				Key:      entry.Key,
				Message:  fmt.Sprintf("duplicate key (first seen on line %d)", prev),
				Severity: SeverityError,
			})
		} else {
			seen[entry.Key] = line
		}

		for _, rule := range l.rules {
			if rule == nil {
				continue
			}
			findings = append(findings, rule(line, entry)...)
		}
	}
	return findings
}

// RuleUpperCaseKeys warns when a key contains lowercase letters.
func RuleUpperCaseKeys(line int, entry parser.Entry) []Finding {
	if entry.Key != strings.ToUpper(entry.Key) {
		return []Finding{{
			Line:     line,
			Key:      entry.Key,
			Message:  "key should be UPPER_CASE",
			Severity: SeverityWarning,
		}}
	}
	return nil
}

// RuleNoDuplicateKeys is handled inside Lint directly; this stub satisfies the Rule type.
func RuleNoDuplicateKeys(_ int, _ parser.Entry) []Finding { return nil }

// RuleNoWhitespaceInValue warns when an unquoted value contains leading/trailing whitespace.
func RuleNoWhitespaceInValue(line int, entry parser.Entry) []Finding {
	if entry.Value != strings.TrimSpace(entry.Value) {
		return []Finding{{
			Line:     line,
			Key:      entry.Key,
			Message:  "value has leading or trailing whitespace",
			Severity: SeverityWarning,
		}}
	}
	return nil
}
