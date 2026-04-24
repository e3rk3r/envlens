// Package sanitize provides tools to strip or redact sensitive entries
// from parsed .env files before they are written to output or shared.
//
// Usage:
//
//	cfg, err := sanitize.LoadConfig("sanitize.json")
//	if err != nil { ... }
//
//	s := cfg.ToSanitizer()
//	clean := s.Apply(entries)
//
// Rules support:
//   - ExactKeys   – drop entries with these exact key names
//   - KeyPrefix   – drop entries whose key starts with the given prefix
//   - KeySuffix   – drop entries whose key ends with the given suffix
//   - RedactKeys  – keep the entry but replace its value with a placeholder
//
// The default placeholder is "***". Override it per-rule via Placeholder.
package sanitize
