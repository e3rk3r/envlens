package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldRule defines validation rules for a single env key.
type FieldRule struct {
	Required bool
	Pattern  string
	Allowed  []string
}

// Schema maps env key names to their validation rules.
type Schema map[string]FieldRule

// Violation represents a single schema validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks a map of env entries against the schema.
// It returns a slice of Violations (empty means valid).
func Validate(schema Schema, env map[string]string) []Violation {
	var violations []Violation

	for key, rule := range schema {
		val, present := env[key]

		if rule.Required && !present {
			violations = append(violations, Violation{Key: key, Message: "required key is missing"})
			continue
		}

		if !present {
			continue
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{Key: key, Message: fmt.Sprintf("invalid pattern in schema: %v", err)})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern)})
			}
		}

		if len(rule.Allowed) > 0 {
			if !containsString(rule.Allowed, val) {
				violations = append(violations, Violation{Key: key, Message: fmt.Sprintf("value %q not in allowed set [%s]", val, strings.Join(rule.Allowed, ", "))})
			}
		}
	}

	return violations
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
