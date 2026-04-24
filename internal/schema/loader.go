package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// schemaFile is the JSON representation used for loading from disk.
type schemaFile struct {
	Fields map[string]struct {
		Required bool     `json:"required"`
		Pattern  string   `json:"pattern"`
		Allowed  []string `json:"allowed"`
	} `json:"fields"`
}

// LoadFromFile reads a JSON schema definition from the given path
// and returns a Schema ready for use with Validate.
//
// Example JSON:
//
//	{
//	  "fields": {
//	    "APP_ENV":  { "required": true, "allowed": ["development","production"] },
//	    "PORT":     { "required": true, "pattern": "^\\d+$" }
//	  }
//	}
func LoadFromFile(path string) (Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: reading file %q: %w", path, err)
	}
	return LoadFromBytes(data)
}

// LoadFromBytes parses a JSON-encoded schema definition.
func LoadFromBytes(data []byte) (Schema, error) {
	var sf schemaFile
	if err := json.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("schema: parsing JSON: %w", err)
	}

	schema := make(Schema, len(sf.Fields))
	for key, f := range sf.Fields {
		schema[key] = FieldRule{
			Required: f.Required,
			Pattern:  f.Pattern,
			Allowed:  f.Allowed,
		}
	}
	return schema, nil
}
