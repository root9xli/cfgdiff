package audit_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/cfgdiff/internal/audit"
	"github.com/cfgdiff/internal/diff"
)

func TestPrintEntries_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	audit.PrintEntries(&buf, nil)
	if !strings.Contains(buf.String(), "No audit entries") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestPrintEntries_WithEntries(t *testing.T) {
	changes := []diff.Change{
		{Key: "key1", Type: diff.Added, NewValue: "val"},
		{Key: "key2", Type: diff.Removed, OldValue: "old"},
	}
	entries := []audit.Entry{
		{
			Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
			FileA:     "prod.yaml",
			FileB:     "staging.yaml",
			Changes:   changes,
			Summary: audit.Summary{
				Added: 1, Removed: 1, Modified: 0, Total: 2,
			},
		},
	}

	var buf bytes.Buffer
	audit.PrintEntries(&buf, entries)
	out := buf.String()

	if !strings.Contains(out, "prod.yaml") {
		t.Errorf("expected prod.yaml in output")
	}
	if !strings.Contains(out, "staging.yaml") {
		t.Errorf("expected staging.yaml in output")
	}
	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected timestamp in output")
	}
	if !strings.Contains(out, "TIMESTAMP") {
		t.Errorf("expected header row in output")
	}
}

func TestPrintEntries_MultipleRows(t *testing.T) {
	entries := []audit.Entry{
		{FileA: "a.env", FileB: "b.env", Summary: audit.Summary{Total: 1}},
		{FileA: "c.json", FileB: "d.json", Summary: audit.Summary{Total: 3}},
	}
	var buf bytes.Buffer
	audit.PrintEntries(&buf, entries)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + 2 data rows = 4
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d: %v", len(lines), lines)
	}
}
