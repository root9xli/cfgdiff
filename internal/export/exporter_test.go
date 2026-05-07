package export

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/diff"
)

var sampleChanges = []diff.Change{
	{Key: "db.host", Type: diff.Added, NewValue: "localhost"},
	{Key: "db.port", Type: diff.Modified, OldValue: "5432", NewValue: "5433"},
	{Key: "db.name", Type: diff.Removed, OldValue: "mydb"},
}

func TestNew_ValidFormat(t *testing.T) {
	for _, f := range []Format{FormatCSV, FormatMarkdown, FormatJSON} {
		_, err := New(f, &bytes.Buffer{})
		if err != nil {
			t.Errorf("expected no error for format %s, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := New("xml", &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatCSV, &buf)
	if err := e.Write(sampleChanges); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "key,type,old_value,new_value") {
		t.Error("CSV missing header row")
	}
	if !strings.Contains(out, "db.host") {
		t.Error("CSV missing expected key")
	}
}

func TestWriteMarkdown(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatMarkdown, &buf)
	if err := e.Write(sampleChanges); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| Key |") {
		t.Error("Markdown missing header")
	}
	if !strings.Contains(out, "db.port") {
		t.Error("Markdown missing expected key")
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatJSON, &buf)
	if err := e.Write(sampleChanges); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []diff.Change
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(out) != len(sampleChanges) {
		t.Errorf("expected %d changes, got %d", len(sampleChanges), len(out))
	}
}

func TestWrite_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	e, _ := New(FormatCSV, &buf)
	if err := e.Write([]diff.Change{}); err != nil {
		t.Fatalf("unexpected error on empty changes: %v", err)
	}
	if !strings.Contains(buf.String(), "key") {
		t.Error("expected header even with empty changes")
	}
}
