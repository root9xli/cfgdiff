package snapshot_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cfgdiff/internal/snapshot"
)

func TestPrintList_Empty(t *testing.T) {
	var buf bytes.Buffer
	snapshot.PrintList(&buf, []string{})
	if !strings.Contains(buf.String(), "No snapshots") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrintList_WithLabels(t *testing.T) {
	var buf bytes.Buffer
	snapshot.PrintList(&buf, []string{"prod", "staging"})
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Error("expected prod in output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected staging in output")
	}
	if !strings.Contains(out, "LABEL") {
		t.Error("expected header LABEL")
	}
}

func TestPrintEntry(t *testing.T) {
	entry := &snapshot.Entry{
		Label:     "prod",
		File:      "config.yaml",
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Data:      map[string]interface{}{"a": "1", "b": "2"},
	}
	var buf bytes.Buffer
	snapshot.PrintEntry(&buf, entry)
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Error("expected label prod")
	}
	if !strings.Contains(out, "config.yaml") {
		t.Error("expected file config.yaml")
	}
	if !strings.Contains(out, "2") {
		t.Error("expected key count 2")
	}
}
