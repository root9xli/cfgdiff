// Package env provides utilities for resolving environment variable
// substitutions within parsed config maps before diffing.
package env

import (
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Resolver replaces environment variable references in config values.
type Resolver struct {
	lookup func(string) (string, bool)
}

// New returns a Resolver that uses os.LookupEnv by default.
func New() *Resolver {
	return &Resolver{lookup: os.LookupEnv}
}

// NewWithMap returns a Resolver backed by a static map, useful for testing.
func NewWithMap(env map[string]string) *Resolver {
	return &Resolver{
		lookup: func(key string) (string, bool) {
			v, ok := env[key]
			return v, ok
		},
	}
}

// Resolve walks a flat config map and expands environment variable references
// found in string values. Non-string values are left untouched.
func (r *Resolver) Resolve(cfg map[string]any) map[string]any {
	out := make(map[string]any, len(cfg))
	for k, v := range cfg {
		switch s := v.(type) {
		case string:
			out[k] = r.expand(s)
		default:
			out[k] = v
		}
	}
	return out
}

// expand replaces all variable references in s with their resolved values.
// If a variable is not found in the environment the reference is left as-is.
func (r *Resolver) expand(s string) string {
	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		key := strings.TrimPrefix(match, "$")
		key = strings.TrimPrefix(key, "{")
		key = strings.TrimSuffix(key, "}")
		if val, ok := r.lookup(key); ok {
			return val
		}
		return match
	})
}
