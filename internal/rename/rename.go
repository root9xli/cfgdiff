// Package rename provides utilities for tracking and applying key renames
// across config diff changes.
package rename

import (
	"strings"

	"github.com/cfgdiff/internal/diff"
)

// Rule defines a single key rename mapping.
type Rule struct {
	From string
	To   string
}

// Options controls rename behaviour.
type Options struct {
	CaseSensitive bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{CaseSensitive: true}
}

// Renamer applies rename rules to a flat config map and to diff change slices.
type Renamer struct {
	rules []Rule
	opts  Options
}

// New creates a Renamer with the given rules and options.
func New(rules []Rule, opts Options) *Renamer {
	return &Renamer{rules: rules, opts: opts}
}

// ApplyToMap returns a new map with keys renamed according to the rules.
// Original map is not mutated.
func (r *Renamer) ApplyToMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[r.rename(k)] = v
	}
	return out
}

// ApplyToChanges returns a new slice of changes with keys renamed.
func (r *Renamer) ApplyToChanges(changes []diff.Change) []diff.Change {
	out := make([]diff.Change, len(changes))
	for i, c := range changes {
		c.Key = r.rename(c.Key)
		out[i] = c
	}
	return out
}

// rename returns the renamed key, or the original if no rule matches.
func (r *Renamer) rename(key string) string {
	for _, rule := range r.rules {
		from := rule.From
		k := key
		if !r.opts.CaseSensitive {
			from = strings.ToLower(from)
			k = strings.ToLower(key)
		}
		if k == from {
			return rule.To
		}
	}
	return key
}

// LoadRules parses a slice of "from=to" strings into Rule values.
func LoadRules(pairs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			continue
		}
		rules = append(rules, Rule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
