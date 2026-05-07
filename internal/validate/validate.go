package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single validation rule applied to config key-value pairs.
type Rule struct {
	Name    string
	Check   func(key, value string) *Violation
}

// Violation represents a failed validation rule.
type Violation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

// Validator holds a set of rules and runs them against a config map.
type Validator struct {
	rules []Rule
}

// New returns a Validator with the default rule set.
func New() *Validator {
	return &Validator{rules: defaultRules()}
}

// NewWithRules returns a Validator using only the provided rules.
func NewWithRules(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// Run validates all key-value pairs in cfg and returns any violations found.
func (v *Validator) Run(cfg map[string]string) []Violation {
	var violations []Violation
	for key, value := range cfg {
		for _, rule := range v.rules {
			if viol := rule.Check(key, value); viol != nil {
				violations = append(violations, *viol)
			}
		}
	}
	return violations
}

func defaultRules() []Rule {
	return []Rule{
		{
			Name: "no-empty-value",
			Check: func(key, value string) *Violation {
				if strings.TrimSpace(value) == "" {
					return &Violation{Key: key, Value: value, Rule: "no-empty-value", Message: "value must not be empty"}
				}
				return nil
			},
		},
		{
			Name: "no-whitespace-key",
			Check: func(key, value string) *Violation {
				if strings.ContainsAny(key, " \t") {
					return &Violation{Key: key, Value: value, Rule: "no-whitespace-key", Message: "key must not contain whitespace"}
				}
				return nil
			},
		},
		{
			Name: "valid-key-format",
			Check: func(key, value string) *Violation {
				matched, _ := regexp.MatchString(`^[A-Za-z_][A-Za-z0-9_.\-]*$`, key)
				if !matched {
					return &Violation{Key: key, Value: value, Rule: "valid-key-format",
						Message: fmt.Sprintf("key %q does not match expected format", key)}
				}
				return nil
			},
		},
	}
}
