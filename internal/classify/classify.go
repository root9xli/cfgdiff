// Package classify categorises diff changes by severity and domain.
package classify

import (
	"strings"

	"github.com/user/cfgdiff/internal/diff"
)

// Severity represents how critical a change is.
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Result holds a classified change.
type Result struct {
	Change   diff.Change
	Severity Severity
	Domain   string
}

// Options controls classification behaviour.
type Options struct {
	// CriticalPatterns are key substrings that map to critical severity.
	CriticalPatterns []string
	// HighPatterns are key substrings that map to high severity.
	HighPatterns []string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		CriticalPatterns: []string{"secret", "password", "token", "key", "cert"},
		HighPatterns:     []string{"host", "port", "url", "endpoint", "dsn", "db"},
	}
}

// Classifier assigns severity and domain to changes.
type Classifier struct {
	opts Options
}

// New creates a Classifier with the given options.
func New(opts Options) *Classifier {
	return &Classifier{opts: opts}
}

// Apply classifies a slice of changes and returns Results.
func (c *Classifier) Apply(changes []diff.Change) []Result {
	results := make([]Result, 0, len(changes))
	for _, ch := range changes {
		results = append(results, Result{
			Change:   ch,
			Severity: c.severity(ch.Key),
			Domain:   c.domain(ch.Key),
		})
	}
	return results
}

func (c *Classifier) severity(key string) Severity {
	lower := strings.ToLower(key)
	for _, p := range c.opts.CriticalPatterns {
		if strings.Contains(lower, p) {
			return SeverityCritical
		}
	}
	for _, p := range c.opts.HighPatterns {
		if strings.Contains(lower, p) {
			return SeverityHigh
		}
	}
	return SeverityLow
}

func (c *Classifier) domain(key string) string {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) > 1 {
		return parts[0]
	}
	if idx := strings.Index(key, "_"); idx > 0 {
		return strings.ToLower(key[:idx])
	}
	return "general"
}
