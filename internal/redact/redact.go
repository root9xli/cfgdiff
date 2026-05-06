package redact

import (
	"regexp"
	"strings"
)

// DefaultPatterns are common sensitive key patterns to redact.
var DefaultPatterns = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"private_key",
	"auth",
	"credential",
}

// Redactor masks sensitive values in config maps.
type Redactor struct {
	patterns []*regexp.Regexp
	mask     string
}

// New creates a Redactor with the given key patterns and mask string.
// If patterns is nil, DefaultPatterns are used.
func New(patterns []string, mask string) (*Redactor, error) {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	if mask == "" {
		mask = "***REDACTED***"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(`(?i)` + regexp.QuoteMeta(p))
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled, mask: mask}, nil
}

// IsSensitive returns true if the key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, re := range r.patterns {
		if re.MatchString(lower) {
			return true
		}
	}
	return false
}

// Apply returns a copy of the map with sensitive values replaced by the mask.
func (r *Redactor) Apply(data map[string]any) map[string]any {
	result := make(map[string]any, len(data))
	for k, v := range data {
		if r.IsSensitive(k) {
			result[k] = r.mask
		} else {
			result[k] = v
		}
	}
	return result
}
