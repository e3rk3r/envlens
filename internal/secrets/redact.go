package secrets

import "strings"

const redactedPlaceholder = "[REDACTED]"

// Redactor masks secret values in env maps based on detected findings.
type Redactor struct {
	detector *Detector
}

// NewRedactor returns a Redactor backed by the given Detector.
func NewRedactor(d *Detector) *Redactor {
	return &Redactor{detector: d}
}

// Redact returns a copy of entries with secret values replaced by a placeholder.
// The set of redacted keys is also returned for reporting.
func (r *Redactor) Redact(entries map[string]string) (redacted map[string]string, keys []string) {
	findings := r.detector.Scan(entries)
	secretKeys := make(map[string]struct{}, len(findings))
	for _, f := range findings {
		secretKeys[f.Key] = struct{}{}
	}

	redacted = make(map[string]string, len(entries))
	for k, v := range entries {
		if _, found := secretKeys[k]; found {
			redacted[k] = redactedPlaceholder
			keys = append(keys, k)
		} else {
			redacted[k] = v
		}
	}
	return redacted, keys
}

// MaskValue partially masks a secret value, revealing only the first 4 chars.
func MaskValue(v string) string {
	v = strings.TrimSpace(v)
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:4] + strings.Repeat("*", len(v)-4)
}
