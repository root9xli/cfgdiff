package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// FieldRule defines validation constraints for a single config key.
type FieldRule struct {
	Required bool
	Type     string // "string", "int", "bool", "float"
	Pattern  string
	Allowed  []string
}

// Schema holds validation rules keyed by config field name.
type Schema struct {
	rules map[string]FieldRule
}

// Load parses a schema definition map into a Schema.
func Load(defs map[string]FieldRule) *Schema {
	return &Schema{rules: defs}
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks a flat config map against the schema rules.
// Returns a slice of violations; empty means valid.
func (s *Schema) Validate(cfg map[string]string) []Violation {
	var violations []Violation

	for key, rule := range s.rules {
		val, exists := cfg[key]
		if !exists || val == "" {
			if rule.Required {
				violations = append(violations, Violation{Key: key, Message: "required field is missing"})
			}
			continue
		}

		if err := checkType(val, rule.Type); err != nil {
			violations = append(violations, Violation{Key: key, Message: err.Error()})
		}

		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err == nil && !re.MatchString(val) {
				violations = append(violations, Violation{Key: key, Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern)})
			}
		}

		if len(rule.Allowed) > 0 && !contains(rule.Allowed, val) {
			violations = append(violations, Violation{Key: key, Message: fmt.Sprintf("value %q not in allowed set [%s]", val, strings.Join(rule.Allowed, ", "))})
		}
	}

	return violations
}

func checkType(val, typ string) error {
	switch typ {
	case "int":
		if !regexp.MustCompile(`^-?\d+$`).MatchString(val) {
			return fmt.Errorf("expected int, got %q", val)
		}
	case "float":
		if !regexp.MustCompile(`^-?\d+(\.\d+)?$`).MatchString(val) {
			return fmt.Errorf("expected float, got %q", val)
		}
	case "bool":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" {
			return fmt.Errorf("expected bool, got %q", val)
		}
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
