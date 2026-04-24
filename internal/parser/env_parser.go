package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single key-value pair parsed from a .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Entries []Entry
	Raw     map[string]string
}

// Parse reads from r and returns a parsed EnvFile.
// It supports KEY=VALUE pairs, inline comments, and skips blank lines.
func Parse(r io.Reader) (*EnvFile, error) {
	ef := &EnvFile{
		Raw: make(map[string]string),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("line %d: invalid format (missing '='): %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		rest := line[idx+1:]

		var value, comment string
		value, comment = splitValueComment(rest)
		value = unquote(strings.TrimSpace(value))

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}

		entry := Entry{
			Key:     key,
			Value:   value,
			Comment: comment,
			Line:    lineNum,
		}
		ef.Entries = append(ef.Entries, entry)
		ef.Raw[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return ef, nil
}

// splitValueComment separates a raw value string from an inline comment.
func splitValueComment(s string) (value, comment string) {
	if idx := strings.Index(s, " #"); idx >= 0 {
		return s[:idx], strings.TrimSpace(s[idx+2:])
	}
	return s, ""
}

// unquote strips surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
