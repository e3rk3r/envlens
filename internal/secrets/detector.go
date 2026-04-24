package secrets

import (
	"regexp"
	"strings"
)

// Finding represents a detected secret in an env entry.
type Finding struct {
	Key     string
	Value   string
	RuleID  string
	Message string
}

// Rule defines a pattern-based secret detection rule.
type Rule struct {
	ID      string
	Message string
	KeyRe   *regexp.Regexp
	ValRe   *regexp.Regexp
}

var defaultRules = []Rule{
	{
		ID:      "high-entropy-secret",
		Message: "Value looks like a high-entropy secret",
		KeyRe:   regexp.MustCompile(`(?i)(secret|token|password|passwd|pwd|api_key|apikey|private_key)`),
		ValRe:   regexp.MustCompile(`[A-Za-z0-9+/]{20,}`),
	},
	{
		ID:      "aws-access-key",
		Message: "Possible AWS access key ID",
		KeyRe:   nil,
		ValRe:   regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	},
	{
		ID:      "generic-bearer-token",
		Message: "Possible bearer/JWT token",
		KeyRe:   nil,
		ValRe:   regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`),
	},
}

// Detector scans env entries for secrets using a set of rules.
type Detector struct {
	Rules []Rule
}

// NewDetector returns a Detector loaded with the default rule set.
func NewDetector() *Detector {
	return &Detector{Rules: defaultRules}
}

// Scan checks a map of key→value pairs and returns all findings.
func (d *Detector) Scan(entries map[string]string) []Finding {
	var findings []Finding
	for k, v := range entries {
		if strings.TrimSpace(v) == "" {
			continue
		}
		for _, rule := range d.Rules {
			keyMatch := rule.KeyRe == nil || rule.KeyRe.MatchString(k)
			valMatch := rule.ValRe == nil || rule.ValRe.MatchString(v)
			if keyMatch && valMatch {
				findings = append(findings, Finding{
					Key:     k,
					Value:   v,
					RuleID:  rule.ID,
					Message: rule.Message,
				})
				break // one finding per key is enough
			}
		}
	}
	return findings
}
