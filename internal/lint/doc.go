// Package lint implements configurable linting rules for .env files.
//
// It checks entries parsed by the parser package against a set of rules
// and returns structured findings with severity levels (error, warning, info).
//
// # Built-in Rules
//
//   - RuleUpperCaseKeys: warns when a key is not fully upper-case.
//   - RuleNoDuplicateKeys: reports keys that appear more than once (always active).
//   - RuleNoWhitespaceInValue: warns when a value has leading or trailing whitespace.
//
// # Usage
//
//	l := lint.New()                 // default rules
//	findings := l.Lint(entries)
//	for _, f := range findings {
//	    fmt.Printf("[%s] line %d %s: %s\n", f.Severity, f.Line, f.Key, f.Message)
//	}
//
// Rules can also be driven by a JSON config file via LoadConfig / Config.ToLinter.
package lint
