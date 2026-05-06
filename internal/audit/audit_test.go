package audit_test

import (
	"os"
	"testing"
	"time"

	"github.com/cfgdiff/internal/audit"
	"github.com/cfgdiff/internal/diff"
)

func sampleChanges() []diff.Change {
	return []diff.Change{
		{Key: "host", Type: diff.Added, NewValue: "localhost"},
		{Key: "port", Type: diff.Modified, OldValue: "8080", NewValue: "9090"},
		{Key: "debug", Type: diff.Removed, OldValue: "true"},
	}
}

func TestRecord_CreatesEntry(t *testing.T) {
	tmp, err := os.CreateTemp("", "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	logger := audit.NewLogger(tmp.Name())
	if err := logger.Record("a.yaml", "b.yaml", sampleChanges()); err != nil {
		t.Fatalf("Record returned error: %v", err)
	}

	entries, err := audit.ReadLog(tmp.Name(), nil)
	if err != nil {
		t.Fatalf("ReadLog error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.FileA != "a.yaml" || e.FileB != "b.yaml" {
		t.Errorf("unexpected files: %s %s", e.FileA, e.FileB)
	}
	if e.Summary.Added != 1 || e.Summary.Removed != 1 || e.Summary.Modified != 1 {
		t.Errorf("unexpected summary: %+v", e.Summary)
	}
	if e.Summary.Total != 3 {
		t.Errorf("expected total 3, got %d", e.Summary.Total)
	}
}

func TestReadLog_EmptyFile(t *testing.T) {
	tmp, _ := os.CreateTemp("", "audit-*.jsonl")
	tmp.Close()
	defer os.Remove(tmp.Name())

	entries, err := audit.ReadLog(tmp.Name(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestReadLog_FilterBySince(t *testing.T) {
	tmp, _ := os.CreateTemp("", "audit-*.jsonl")
	tmp.Close()
	defer os.Remove(tmp.Name())

	logger := audit.NewLogger(tmp.Name())
	logger.Record("old.yaml", "new.yaml", sampleChanges())

	future := time.Now().Add(time.Hour)
	f := &audit.Filter{Since: &future}
	entries, err := audit.ReadLog(tmp.Name(), f)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries after future filter, got %d", len(entries))
	}
}

func TestReadLog_MissingFile(t *testing.T) {
	entries, err := audit.ReadLog("/nonexistent/path/audit.jsonl", nil)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil slice for missing file")
	}
}
