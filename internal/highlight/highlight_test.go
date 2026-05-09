package highlight

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/diff"
)

func sampleChanges() []diff.Change {
	return []diff.Change{
		{Key: "APP_ENV", Type: diff.Added, NewValue: "production"},
		{Key: "DB_HOST", Type: diff.Removed, OldValue: "localhost"},
		{Key: "LOG_LEVEL", Type: diff.Modified, OldValue: "debug", NewValue: "warn"},
	}
}

func TestLine_Added_NoColor(t *testing.T) {
	h := New(Options{Enabled: false})
	line := h.Line(diff.Change{Key: "X", Type: diff.Added, NewValue: "1"})
	if !strings.Contains(line, "+ X = 1") {
		t.Errorf("unexpected line: %q", line)
	}
}

func TestLine_Removed_NoColor(t *testing.T) {
	h := New(Options{Enabled: false})
	line := h.Line(diff.Change{Key: "Y", Type: diff.Removed, OldValue: "old"})
	if !strings.Contains(line, "- Y = old") {
		t.Errorf("unexpected line: %q", line)
	}
}

func TestLine_Modified_NoColor(t *testing.T) {
	h := New(Options{Enabled: false})
	line := h.Line(diff.Change{Key: "Z", Type: diff.Modified, OldValue: "a", NewValue: "b"})
	if !strings.Contains(line, "~ Z") || !strings.Contains(line, "a") || !strings.Contains(line, "b") {
		t.Errorf("unexpected line: %q", line)
	}
}

func TestLine_Added_WithColor(t *testing.T) {
	h := New(DefaultOptions())
	line := h.Line(diff.Change{Key: "K", Type: diff.Added, NewValue: "v"})
	if !strings.Contains(line, colorGreen) {
		t.Errorf("expected green color code in: %q", line)
	}
}

func TestApply_ReturnsAllLines(t *testing.T) {
	h := New(Options{Enabled: false})
	lines := h.Apply(sampleChanges())
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestHeader_NoColor(t *testing.T) {
	h := New(Options{Enabled: false})
	out := h.Header("Summary")
	if out != "Summary" {
		t.Errorf("expected plain text, got %q", out)
	}
}

func TestHeader_WithColor(t *testing.T) {
	h := New(DefaultOptions())
	out := h.Header("Summary")
	if !strings.Contains(out, colorCyan) {
		t.Errorf("expected cyan color code in header: %q", out)
	}
}

func TestFprint_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, []diff.Change{}, Options{Enabled: false})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %q", buf.String())
	}
}

func TestFprint_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, sampleChanges(), Options{Enabled: false})
	out := buf.String()
	if !strings.Contains(out, "3 change(s)") {
		t.Errorf("expected change count in header, got: %q", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in output, got: %q", out)
	}
}
