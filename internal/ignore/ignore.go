package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Rules holds a set of key patterns to ignore during diff.
type Rules struct {
	patterns []string
}

// LoadFile reads ignore patterns from a file (one pattern per line).
// Lines starting with '#' and empty lines are skipped.
func LoadFile(path string) (*Rules, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Rules{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var patterns []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Rules{patterns: patterns}, nil
}

// NewRules creates Rules from a slice of pattern strings.
func NewRules(patterns []string) *Rules {
	return &Rules{patterns: patterns}
}

// Match reports whether the given key matches any ignore pattern.
// Patterns support shell-style wildcards via filepath.Match.
func (r *Rules) Match(key string) bool {
	for _, p := range r.patterns {
		matched, err := filepath.Match(p, key)
		if err == nil && matched {
			return true
		}
		// Also support prefix matching for nested keys (e.g. "secret.*")
		if strings.HasSuffix(p, ".*") {
			prefix := strings.TrimSuffix(p, ".*")
			if strings.HasPrefix(key, prefix+".") {
				return true
			}
		}
	}
	return false
}

// FilterKeys returns only the keys from the provided slice that are NOT ignored.
func (r *Rules) FilterKeys(keys []string) []string {
	if len(r.patterns) == 0 {
		return keys
	}
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if !r.Match(k) {
			out = append(out, k)
		}
	}
	return out
}

// Patterns returns a copy of the ignore patterns held by this Rules instance.
func (r *Rules) Patterns() []string {
	if len(r.patterns) == 0 {
		return nil
	}
	out := make([]string, len(r.patterns))
	copy(out, r.patterns)
	return out
}
