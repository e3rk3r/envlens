package compare_test

import (
	"testing"

	"github.com/user/envlens/internal/compare"
	"github.com/user/envlens/internal/parser"
)

func entries(pairs ...string) []parser.Entry {
	var out []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompare_AllEnvsHaveSameKeys(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("DB_HOST", "localhost", "PORT", "5432"),
		"prod": entries("DB_HOST", "db.prod.example.com", "PORT", "5432"),
	}
	m := compare.Compare(envs)
	if len(m.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(m.Keys))
	}
	if len(m.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", m.Missing)
	}
	if m.Cells["DB_HOST"]["prod"] != "db.prod.example.com" {
		t.Errorf("unexpected value for DB_HOST in prod")
	}
}

func TestCompare_MissingKeyInOneEnv(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":  entries("DB_HOST", "localhost", "SECRET", "abc"),
		"prod": entries("DB_HOST", "prod-host"),
	}
	m := compare.Compare(envs)
	if len(m.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(m.Keys))
	}
	missingEnvs, ok := m.Missing["SECRET"]
	if !ok {
		t.Fatal("expected SECRET to be missing in some env")
	}
	if len(missingEnvs) != 1 || missingEnvs[0] != "prod" {
		t.Errorf("expected SECRET missing in prod, got %v", missingEnvs)
	}
}

func TestCompare_ThreeEnvironments(t *testing.T) {
	envs := map[string][]parser.Entry{
		"dev":     entries("A", "1", "B", "2"),
		"staging": entries("A", "1"),
		"prod":    entries("A", "1", "B", "2", "C", "3"),
	}
	m := compare.Compare(envs)
	if len(m.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(m.Keys))
	}
	if len(m.Envs) != 3 {
		t.Fatalf("expected 3 envs, got %d", len(m.Envs))
	}
	if len(m.Missing["B"]) != 1 {
		t.Errorf("expected B missing in 1 env, got %v", m.Missing["B"])
	}
	if len(m.Missing["C"]) != 2 {
		t.Errorf("expected C missing in 2 envs, got %v", m.Missing["C"])
	}
}

func TestCompare_EmptyEnvs(t *testing.T) {
	m := compare.Compare(map[string][]parser.Entry{})
	if len(m.Keys) != 0 {
		t.Errorf("expected no keys for empty input")
	}
}
