// Package pivot provides utilities for transposing config diff results
// into a key-centric view across multiple named environments.
package pivot

import (
	"sort"

	"github.com/user/cfgdiff/internal/diff"
)

// Row represents a single config key and its values across environments.
type Row struct {
	Key    string
	Values map[string]string // env name -> value ("<missing>" if absent)
}

// Table is the result of pivoting multiple environment diffs.
type Table struct {
	Envs []string // ordered environment names
	Rows []Row
}

// missing is the sentinel used when a key is absent in an environment.
const missing = "<missing>"

// Build constructs a pivot Table from a map of environment name to flat
// config map. Each key that appears in any environment becomes a Row.
func Build(envMaps map[string]map[string]string) Table {
	// Collect ordered env names for deterministic output.
	envs := make([]string, 0, len(envMaps))
	for e := range envMaps {
		envs = append(envs, e)
	}
	sort.Strings(envs)

	// Collect all unique keys.
	keySet := map[string]struct{}{}
	for _, m := range envMaps {
		for k := range m {
			keySet[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]Row, 0, len(keys))
	for _, k := range keys {
		vals := make(map[string]string, len(envs))
		for _, e := range envs {
			if v, ok := envMaps[e][k]; ok {
				vals[e] = v
			} else {
				vals[e] = missing
			}
		}
		rows = append(rows, Row{Key: k, Values: vals})
	}

	return Table{Envs: envs, Rows: rows}
}

// Divergent returns only the rows where at least one environment differs
// from the others.
func (t Table) Divergent() []Row {
	var out []Row
	for _, row := range t.Rows {
		if isDivergent(row.Values) {
			out = append(out, row)
		}
	}
	return out
}

func isDivergent(vals map[string]string) bool {
	var first string
	set := false
	for _, v := range vals {
		if !set {
			first = v
			set = true
			continue
		}
		if v != first {
			return true
		}
	}
	return false
}

// FromChanges builds an envMaps input from a named slice of diff.Change lists,
// reconstructing per-environment values from Added/Removed/Modified records.
func FromChanges(base map[string]string, named map[string][]diff.Change) map[string]map[string]string {
	result := map[string]map[string]string{}
	for env, changes := range named {
		m := make(map[string]string, len(base))
		for k, v := range base {
			m[k] = v
		}
		for _, c := range changes {
			switch c.Type {
			case diff.Added:
				m[c.Key] = c.NewValue
			case diff.Removed:
				delete(m, c.Key)
			case diff.Modified:
				m[c.Key] = c.NewValue
			}
		}
		result[env] = m
	}
	return result
}
