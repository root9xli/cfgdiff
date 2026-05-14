package pivot_test

import (
	"testing"

	"github.com/user/cfgdiff/internal/diff"
	"github.com/user/cfgdiff/internal/pivot"
)

func TestBuild_AllEnvsPresent(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"host": "prod.example.com", "port": "443"},
		"staging": {"host": "staging.example.com", "port": "8443"},
	}
	table := pivot.Build(envs)

	if len(table.Envs) != 2 {
		t.Fatalf("expected 2 envs, got %d", len(table.Envs))
	}
	if len(table.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table.Rows))
	}
	// Rows should be sorted by key.
	if table.Rows[0].Key != "host" || table.Rows[1].Key != "port" {
		t.Errorf("unexpected row order: %v", []string{table.Rows[0].Key, table.Rows[1].Key})
	}
}

func TestBuild_MissingKeyInEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"host": "prod.example.com", "debug": "false"},
		"staging": {"host": "staging.example.com"},
	}
	table := pivot.Build(envs)

	for _, row := range table.Rows {
		if row.Key == "debug" {
			if row.Values["staging"] != "<missing>" {
				t.Errorf("expected <missing> for staging debug, got %q", row.Values["staging"])
			}
			return
		}
	}
	t.Error("debug row not found")
}

func TestDivergent_FiltersUniform(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"host": "same", "port": "443"},
		"staging": {"host": "same", "port": "8443"},
	}
	table := pivot.Build(envs)
	div := table.Divergent()

	if len(div) != 1 {
		t.Fatalf("expected 1 divergent row, got %d", len(div))
	}
	if div[0].Key != "port" {
		t.Errorf("expected divergent key 'port', got %q", div[0].Key)
	}
}

func TestDivergent_AllUniform(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"key": "val"},
		"staging": {"key": "val"},
	}
	table := pivot.Build(envs)
	if len(table.Divergent()) != 0 {
		t.Error("expected no divergent rows")
	}
}

func TestFromChanges_AppliesModified(t *testing.T) {
	base := map[string]string{"host": "base.example.com", "port": "80"}
	changes := map[string][]diff.Change{
		"prod": {
			{Key: "host", Type: diff.Modified, OldValue: "base.example.com", NewValue: "prod.example.com"},
		},
	}
	envMaps := pivot.FromChanges(base, changes)
	if envMaps["prod"]["host"] != "prod.example.com" {
		t.Errorf("expected prod host to be updated, got %q", envMaps["prod"]["host"])
	}
	if envMaps["prod"]["port"] != "80" {
		t.Errorf("expected prod port unchanged, got %q", envMaps["prod"]["port"])
	}
}

func TestFromChanges_AppliesAddedAndRemoved(t *testing.T) {
	base := map[string]string{"host": "base.example.com", "legacy": "old"}
	changes := map[string][]diff.Change{
		"staging": {
			{Key: "newkey", Type: diff.Added, NewValue: "newval"},
			{Key: "legacy", Type: diff.Removed, OldValue: "old"},
		},
	}
	envMaps := pivot.FromChanges(base, changes)
	if _, ok := envMaps["staging"]["legacy"]; ok {
		t.Error("expected legacy key to be removed")
	}
	if envMaps["staging"]["newkey"] != "newval" {
		t.Errorf("expected newkey=newval, got %q", envMaps["staging"]["newkey"])
	}
}
