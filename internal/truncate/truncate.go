// Package truncate provides utilities for truncating long config values
// in diff output to improve readability.
package truncate

import "fmt"

const DefaultMaxLen = 80

// Options controls truncation behaviour.
type Options struct {
	MaxLen  int
	Suffix  string
	Enabled bool
}

// DefaultOptions returns sensible truncation defaults.
func DefaultOptions() Options {
	return Options{
		MaxLen:  DefaultMaxLen,
		Suffix:  "...",
		Enabled: true,
	}
}

// Truncator applies truncation to config values.
type Truncator struct {
	opts Options
}

// New creates a Truncator with the given options.
func New(opts Options) *Truncator {
	if opts.MaxLen <= 0 {
		opts.MaxLen = DefaultMaxLen
	}
	if opts.Suffix == "" {
		opts.Suffix = "..."
	}
	return &Truncator{opts: opts}
}

// Value truncates a single string value if it exceeds MaxLen.
func (t *Truncator) Value(s string) string {
	if !t.opts.Enabled {
		return s
	}
	if len(s) <= t.opts.MaxLen {
		return s
	}
	cutoff := t.opts.MaxLen - len(t.opts.Suffix)
	if cutoff < 0 {
		cutoff = 0
	}
	return s[:cutoff] + t.opts.Suffix
}

// Apply truncates all values in the provided map, returning a new map.
// Keys are never truncated.
func (t *Truncator) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = t.Value(v)
	}
	return out
}

// Summary returns a one-line description of how many values were truncated.
func (t *Truncator) Summary(original, truncated map[string]string) string {
	count := 0
	for k, v := range original {
		if truncated[k] != v {
			count++
		}
	}
	if count == 0 {
		return "no values truncated"
	}
	return fmt.Sprintf("%d value(s) truncated to %d chars", count, t.opts.MaxLen)
}
