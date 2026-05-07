package schema

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintViolations_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	PrintViolations(&buf, nil)
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected no-violations message, got: %s", buf.String())
	}
}

func TestPrintViolations_WithViolations(t *testing.T) {
	var buf bytes.Buffer
	v := []Violation{
		{Key: "PORT", Message: "expected int, got \"abc\""},
		{Key: "DB_URL", Message: "required field is missing"},
	}
	PrintViolations(&buf, v)
	out := buf.String()

	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output")
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output")
	}
	if !strings.Contains(out, "2 violation(s)") {
		t.Errorf("expected violation count in output, got: %s", out)
	}
}

func TestPrintViolations_HeaderPresent(t *testing.T) {
	var buf bytes.Buffer
	PrintViolations(&buf, []Violation{{Key: "X", Message: "bad"}})
	out := buf.String()
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "MESSAGE") {
		t.Errorf("expected table header in output, got: %s", out)
	}
}
