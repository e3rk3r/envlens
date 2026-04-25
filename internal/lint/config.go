package lint

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds lint rule toggles loaded from a JSON config file.
type Config struct {
	RequireUpperCase    bool `json:"require_upper_case"`
	ForbidDuplicateKeys bool `json:"forbid_duplicate_keys"`
	ForbidWhitespace    bool `json:"forbid_whitespace"`
}

// DefaultConfig returns a Config with all rules enabled.
func DefaultConfig() Config {
	return Config{
		RequireUpperCase:    true,
		ForbidDuplicateKeys: true,
		ForbidWhitespace:    true,
	}
}

// LoadConfig reads a JSON lint config from the given path.
func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("lint: read config: %w", err)
	}
	return ParseConfig(data)
}

// ParseConfig decodes JSON bytes into a Config.
func ParseConfig(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("lint: parse config: %w", err)
	}
	return cfg, nil
}

// ToLinter builds a Linter from the config, enabling only the selected rules.
func (c Config) ToLinter() *Linter {
	var rules []Rule
	if c.RequireUpperCase {
		rules = append(rules, RuleUpperCaseKeys)
	}
	if c.ForbidWhitespace {
		rules = append(rules, RuleNoWhitespaceInValue)
	}
	// Duplicate key detection is always handled inside Lint();
	// we include the no-op stub only when the rule is active so
	// callers can inspect the rule slice length if needed.
	if c.ForbidDuplicateKeys {
		rules = append(rules, RuleNoDuplicateKeys)
	}
	return WithRules(rules...)
}
