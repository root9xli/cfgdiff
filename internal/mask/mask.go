package mask

import (
	"regexp"
	"strings"
)

// Options controls masking behaviour.
type Options struct {
	Char      string // replacement character, default "*"
	ShowFirst int    // number of leading chars to reveal
	ShowLast  int    // number of trailing chars to reveal
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Char:      "*",
		ShowFirst: 0,
		ShowLast:  0,
	}
}

// Masker applies value masking to config maps.
type Masker struct {
	opts     Options
	patterns []*regexp.Regexp
}

var defaultPatterns = []string{
	`(?i)password`,
	`(?i)secret`,
	`(?i)token`,
	`(?i)api[_-]?key`,
	`(?i)private[_-]?key`,
	`(?i)auth`,
}

// New creates a Masker with default sensitive-key patterns.
func New(opts Options) *Masker {
	return NewWithPatterns(opts, defaultPatterns)
}

// NewWithPatterns creates a Masker with custom key patterns.
func NewWithPatterns(opts Options, patterns []string) *Masker {
	if opts.Char == "" {
		opts.Char = "*"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiled = append(compiled, re)
		}
	}
	return &Masker{opts: opts, patterns: compiled}
}

// IsSensitive reports whether a key matches any pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, re := range m.patterns {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskValue masks a single string value according to options.
func (m *Masker) MaskValue(value string) string {
	n := len(value)
	if n == 0 {
		return value
	}
	show := m.opts.ShowFirst + m.opts.ShowLast
	if show >= n {
		return strings.Repeat(m.opts.Char, n)
	}
	prefix := value[:m.opts.ShowFirst]
	suffix := ""
	if m.opts.ShowLast > 0 {
		suffix = value[n-m.opts.ShowLast:]
	}
	midLen := n - m.opts.ShowFirst - m.opts.ShowLast
	return prefix + strings.Repeat(m.opts.Char, midLen) + suffix
}

// Apply returns a copy of the map with sensitive values masked.
func (m *Masker) Apply(data map[string]string) map[string]string {
	out := make(map[string]string, len(data))
	for k, v := range data {
		if m.IsSensitive(k) {
			out[k] = m.MaskValue(v)
		} else {
			out[k] = v
		}
	}
	return out
}
