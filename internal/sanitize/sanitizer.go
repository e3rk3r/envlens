package sanitize

import (
	"strings"

	"github.com/envlens/internal/parser"
)

// Rule defines a sanitization rule for env entries.
type Rule struct {
	// KeyPrefix removes entries whose keys start with this prefix.
	KeyPrefix string
	// KeySuffix removes entries whose keys end with this suffix.
	KeySuffix string
	// ExactKeys is a set of exact key names to remove.
	ExactKeys []string
	// RedactKeys replaces the value of matching keys with a placeholder.
	RedactKeys []string
	// Placeholder is used when redacting values. Defaults to "***".
	Placeholder string
}

// Sanitizer applies sanitization rules to a slice of env entries.
type Sanitizer struct {
	rules       []Rule
	placeholder string
}

// NewSanitizer creates a Sanitizer with the given rules.
func NewSanitizer(rules []Rule) *Sanitizer {
	ph := "***"
	for _, r := range rules {
		if r.Placeholder != "" {
			ph = r.Placeholder
			break
		}
	}
	return &Sanitizer{rules: rules, placeholder: ph}
}

// Apply runs all rules against entries and returns a sanitized copy.
func (s *Sanitizer) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if s.shouldDrop(e.Key) {
			continue
		}
		if s.shouldRedact(e.Key) {
			e.Value = s.placeholder
		}
		result = append(result, e)
	}
	return result
}

func (s *Sanitizer) shouldDrop(key string) bool {
	for _, r := range s.rules {
		if r.KeyPrefix != "" && strings.HasPrefix(key, r.KeyPrefix) {
			return true
		}
		if r.KeySuffix != "" && strings.HasSuffix(key, r.KeySuffix) {
			return true
		}
		for _, k := range r.ExactKeys {
			if key == k {
				return true
			}
		}
	}
	return false
}

func (s *Sanitizer) shouldRedact(key string) bool {
	for _, r := range s.rules {
		for _, k := range r.RedactKeys {
			if key == k {
				return true
			}
		}
	}
	return false
}
