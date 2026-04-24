package sanitize

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds serializable sanitization configuration.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON sanitizer config from the given file path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("sanitize: read config %q: %w", path, err)
	}
	return ParseConfig(data)
}

// ParseConfig parses sanitizer config from raw JSON bytes.
func ParseConfig(data []byte) (*Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("sanitize: parse config: %w", err)
	}
	return &cfg, nil
}

// ToSanitizer converts a Config into a ready-to-use Sanitizer.
func (c *Config) ToSanitizer() *Sanitizer {
	return NewSanitizer(c.Rules)
}
