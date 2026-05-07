// Package filter provides key-based filtering for config diff results.
package filter

import (
	"strings"

	"github.com/yourorg/cfgdiff/internal/diff"
)

// Options holds filtering criteria for diff changes.
type Options struct {
	// Prefix restricts results to keys starting with the given prefix.
	Prefix string
	// Contains restricts results to keys containing the given substring.
	Contains string
	// Types restricts results to specific change types (added, removed, modified).
	Types []string
}

// Filter applies Options to a slice of diff.Change and returns matching entries.
type Filter struct {
	opts Options
	typeSet map[string]struct{}
}

// New creates a Filter from the given Options.
func New(opts Options) *Filter {
	typeSet := make(map[string]struct{}, len(opts.Types))
	for _, t := range opts.Types {
		typeSet[strings.ToLower(t)] = struct{}{}
	}
	return &Filter{opts: opts, typeSet: typeSet}
}

// Apply returns only the changes that match all active filter criteria.
func (f *Filter) Apply(changes []diff.Change) []diff.Change {
	var out []diff.Change
	for _, c := range changes {
		if !f.matchPrefix(c.Key) {
			continue
		}
		if !f.matchContains(c.Key) {
			continue
		}
		if !f.matchType(c.Type) {
			continue
		}
		out = append(out, c)
	}
	return out
}

func (f *Filter) matchPrefix(key string) bool {
	if f.opts.Prefix == "" {
		return true
	}
	return strings.HasPrefix(key, f.opts.Prefix)
}

func (f *Filter) matchContains(key string) bool {
	if f.opts.Contains == "" {
		return true
	}
	return strings.Contains(key, f.opts.Contains)
}

func (f *Filter) matchType(changeType string) bool {
	if len(f.typeSet) == 0 {
		return true
	}
	_, ok := f.typeSet[strings.ToLower(changeType)]
	return ok
}
