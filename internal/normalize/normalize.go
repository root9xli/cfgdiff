// Package normalize provides key normalization utilities for config maps.
// It supports trimming whitespace, lowercasing keys, and removing common prefixes.
package normalize

import (
	"strings"
)

// Options controls how normalization is applied.
type Options struct {
	LowercaseKeys  bool
	TrimSpace      bool
	StripPrefix    string
}

// DefaultOptions returns a sensible default set of normalization options.
func DefaultOptions() Options {
	return Options{
		LowercaseKeys: true,
		TrimSpace:     true,
		StripPrefix:   "",
	}
}

// Apply normalizes all keys and values in the given config map according to opts.
// It returns a new map and does not mutate the original.
func Apply(cfg map[string]string, opts Options) map[string]string {
	result := make(map[string]string, len(cfg))
	for k, v := range cfg {
		nk := normalizeKey(k, opts)
		nv := normalizeValue(v, opts)
		result[nk] = nv
	}
	return result
}

// NormalizeKey applies normalization rules to a single key.
func NormalizeKey(key string, opts Options) string {
	return normalizeKey(key, opts)
}

func normalizeKey(key string, opts Options) string {
	if opts.TrimSpace {
		key = strings.TrimSpace(key)
	}
	if opts.StripPrefix != "" {
		key = strings.TrimPrefix(key, opts.StripPrefix)
	}
	if opts.LowercaseKeys {
		key = strings.ToLower(key)
	}
	return key
}

func normalizeValue(val string, opts Options) string {
	if opts.TrimSpace {
		return strings.TrimSpace(val)
	}
	return val
}
