// Package flatten converts nested config maps into dot-notation flat maps.
package flatten

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls flatten behaviour.
type Options struct {
	Separator string // key separator, defaults to "."
	Prefix    string // optional prefix prepended to all keys
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{Separator: "."}
}

// Flatten converts a potentially nested map[string]any into a flat
// map[string]string using dot-notation keys.
//
// Example:
//
//	{"db": {"host": "localhost", "port": 5432}} →
//	{"db.host": "localhost", "db.port": "5432"}
func Flatten(input map[string]any, opts Options) map[string]string {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	out := make(map[string]string)
	flatten(input, opts.Prefix, opts.Separator, out)
	return out
}

func flatten(node map[string]any, prefix, sep string, out map[string]string) {
	for k, v := range node {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch val := v.(type) {
		case map[string]any:
			flatten(val, key, sep, out)
		case map[any]any:
			// YAML unmarshals some maps as map[any]any
			converted := make(map[string]any, len(val))
			for mk, mv := range val {
				converted[fmt.Sprintf("%v", mk)] = mv
			}
			flatten(converted, key, sep, out)
		default:
			out[key] = fmt.Sprintf("%v", v)
		}
	}
}

// Keys returns the sorted keys of a flat map.
func Keys(flat map[string]string) []string {
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// HasPrefix reports whether key starts with the given dot-separated prefix.
func HasPrefix(key, prefix, sep string) bool {
	if sep == "" {
		sep = "."
	}
	return key == prefix || strings.HasPrefix(key, prefix+sep)
}
