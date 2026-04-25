// Package compare provides multi-environment .env file comparison,
// allowing users to compare more than two environments at once and
// produce a unified matrix of keys and their values per environment.
package compare

import (
	"sort"

	"github.com/user/envlens/internal/parser"
)

// EnvName is a label for an environment (e.g. "dev", "prod").
type EnvName = string

// Matrix holds the comparison result across multiple environments.
type Matrix struct {
	// Keys is the sorted list of all unique keys found across all envs.
	Keys []string
	// Envs is the ordered list of environment names.
	Envs []EnvName
	// Cells maps key -> envName -> value (empty string if absent).
	Cells map[string]map[EnvName]string
	// Missing maps key -> list of envNames where the key is absent.
	Missing map[string][]EnvName
}

// Compare builds a Matrix from a map of environment name -> parsed entries.
func Compare(envs map[EnvName][]parser.Entry) *Matrix {
	keySet := make(map[string]struct{})
	for _, entries := range envs {
		for _, e := range entries {
			keySet[e.Key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	envNames := make([]EnvName, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	cells := make(map[string]map[EnvName]string, len(keys))
	missing := make(map[string][]EnvName)

	for _, key := range keys {
		cells[key] = make(map[EnvName]string, len(envNames))
	}

	for _, name := range envNames {
		present := make(map[string]string)
		for _, e := range envs[name] {
			present[e.Key] = e.Value
		}
		for _, key := range keys {
			if val, ok := present[key]; ok {
				cells[key][name] = val
			} else {
				cells[key][name] = ""
				missing[key] = append(missing[key], name)
			}
		}
	}

	return &Matrix{
		Keys:    keys,
		Envs:    envNames,
		Cells:   cells,
		Missing: missing,
	}
}
