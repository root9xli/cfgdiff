package lint

import (
	"fmt"
	"strings"
)

// Rule represents a single lint rule applied to config keys/values.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Violation holds a failed lint result.
type Violation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

// Linter runs a set of rules against a flat config map.
type Linter struct {
	rules []Rule
}

// New creates a Linter with the default built-in rules.
func New() *Linter {
	return &Linter{
		rules: defaultRules(),
	}
}

// NewWithRules creates a Linter with custom rules only.
func NewWithRules(rules []Rule) *Linter {
	return &Linter{rules: rules}
}

// Run applies all rules to the provided config map and returns violations.
func (l *Linter) Run(cfg map[string]string) []Violation {
	var violations []Violation
	for key, value := range cfg {
		for _, rule := range l.rules {
			if rule.Check(key, value) {
				violations = append(violations, Violation{
					Key:     key,
					Value:   value,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations
}

func defaultRules() []Rule {
	return []Rule{
		{
			Name:    "empty-value",
			Message: "key has an empty value",
			Check: func(key, value string) bool {
				return strings.TrimSpace(value) == ""
			},
		},
		{
			Name:    "uppercase-key",
			Message: "key should be lowercase",
			Check: func(key, value string) bool {
				return key != strings.ToLower(key)
			},
		},
		{
			Name:    "whitespace-value",
			Message: "value contains leading or trailing whitespace",
			Check: func(key, value string) bool {
				return value != strings.TrimSpace(value) && value != ""
			},
		},
		{
			Name:    "duplicate-dots",
			Message: "key contains consecutive dots",
			Check: func(key, value string) bool {
				return strings.Contains(key, "..")
			},
		},
		{
			Name:    "placeholder-value",
			Message: fmt.Sprintf("value looks like an unset placeholder"),
			Check: func(key, value string) bool {
				v := strings.ToLower(strings.TrimSpace(value))
				return v == "todo" || v == "fixme" || v == "changeme" || v == "<your-value>"
			},
		},
	}
}
