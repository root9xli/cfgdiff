package lint

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_NoViolations(t *testing.T) {
	l := New()
	cfg := map[string]string{
		"app.name":    "cfgdiff",
		"app.version": "1.0.0",
	}
	violations := l.Run(cfg)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestRun_EmptyValue(t *testing.T) {
	l := New()
	cfg := map[string]string{"db.host": ""}
	v := l.Run(cfg)
	if !hasRule(v, "empty-value") {
		t.Error("expected empty-value violation")
	}
}

func TestRun_UppercaseKey(t *testing.T) {
	l := New()
	cfg := map[string]string{"DB_HOST": "localhost"}
	v := l.Run(cfg)
	if !hasRule(v, "uppercase-key") {
		t.Error("expected uppercase-key violation")
	}
}

func TestRun_WhitespaceValue(t *testing.T) {
	l := New()
	cfg := map[string]string{"app.name": "  cfgdiff "}
	v := l.Run(cfg)
	if !hasRule(v, "whitespace-value") {
		t.Error("expected whitespace-value violation")
	}
}

func TestRun_DoubleDotKey(t *testing.T) {
	l := New()
	cfg := map[string]string{"app..name": "cfgdiff"}
	v := l.Run(cfg)
	if !hasRule(v, "duplicate-dots") {
		t.Error("expected duplicate-dots violation")
	}
}

func TestRun_PlaceholderValue(t *testing.T) {
	l := New()
	cfg := map[string]string{"api.key": "changeme"}
	v := l.Run(cfg)
	if !hasRule(v, "placeholder-value") {
		t.Error("expected placeholder-value violation")
	}
}

func TestNewWithRules_CustomOnly(t *testing.T) {
	custom := []Rule{
		{
			Name:    "no-localhost",
			Message: "value must not be localhost in production",
			Check: func(key, value string) bool {
				return value == "localhost"
			},
		},
	}
	l := NewWithRules(custom)
	v := l.Run(map[string]string{"db.host": "localhost"})
	if !hasRule(v, "no-localhost") {
		t.Error("expected no-localhost violation")
	}
}

func TestPrintViolations_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	PrintViolations(&buf, nil)
	if !strings.Contains(buf.String(), "no lint violations") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintViolations_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	v := []Violation{{Key: "DB_HOST", Rule: "uppercase-key", Message: "key should be lowercase"}}
	PrintViolations(&buf, v)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "uppercase-key") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected violation count in output: %s", out)
	}
}

func hasRule(violations []Violation, name string) bool {
	for _, v := range violations {
		if v.Rule == name {
			return true
		}
	}
	return false
}
