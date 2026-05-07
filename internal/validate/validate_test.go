package validate

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoViolations(t *testing.T) {
	v := New()
	cfg := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	got := v.Run(cfg)
	if len(got) != 0 {
		t.Errorf("expected 0 violations, got %d", len(got))
	}
}

func TestRun_EmptyValue(t *testing.T) {
	v := New()
	cfg := map[string]string{"DB_PASS": ""}
	got := v.Run(cfg)
	if len(got) == 0 {
		t.Fatal("expected a violation for empty value")
	}
	if got[0].Rule != "no-empty-value" {
		t.Errorf("unexpected rule: %s", got[0].Rule)
	}
}

func TestRun_WhitespaceKey(t *testing.T) {
	v := New()
	cfg := map[string]string{"bad key": "value"}
	got := v.Run(cfg)
	rules := map[string]bool{}
	for _, viol := range got {
		rules[viol.Rule] = true
	}
	if !rules["no-whitespace-key"] {
		t.Error("expected no-whitespace-key violation")
	}
}

func TestRun_InvalidKeyFormat(t *testing.T) {
	v := New()
	cfg := map[string]string{"123invalid": "value"}
	got := v.Run(cfg)
	rules := map[string]bool{}
	for _, viol := range got {
		rules[viol.Rule] = true
	}
	if !rules["valid-key-format"] {
		t.Error("expected valid-key-format violation")
	}
}

func TestRun_CustomRules(t *testing.T) {
	customRule := Rule{
		Name: "no-localhost",
		Check: func(key, value string) *Violation {
			if value == "localhost" {
				return &Violation{Key: key, Value: value, Rule: "no-localhost", Message: "localhost not allowed in production"}
			}
			return nil
		},
	}
	v := NewWithRules([]Rule{customRule})
	cfg := map[string]string{"HOST": "localhost"}
	got := v.Run(cfg)
	if len(got) != 1 || got[0].Rule != "no-localhost" {
		t.Errorf("expected no-localhost violation, got %+v", got)
	}
}

func TestPrintViolations_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	PrintViolations(&buf, nil)
	if !strings.Contains(buf.String(), "no validation violations") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintViolations_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	viols := []Violation{
		{Key: "FOO", Value: "", Rule: "no-empty-value", Message: "value must not be empty"},
	}
	PrintViolations(&buf, viols)
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Error("expected key FOO in output")
	}
	if !strings.Contains(out, "1 violation(s)") {
		t.Error("expected violation count in output")
	}
}
