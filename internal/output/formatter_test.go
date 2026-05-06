package output

import (
	"bytes"
	"strings"
	"testing"

	"cfgdiff/internal/diff"
)

func changes() []diff.Change {
	return []diff.Change{
		{Key: "db.host", Type: diff.Added, NewValue: "localhost"},
		{Key: "db.port", Type: diff.Removed, OldValue: "5432"},
		{Key: "app.env", Type: diff.Modified, OldValue: "dev", NewValue: "prod"},
	}
}

func TestWriteText_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText, true)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestWriteText_WithChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText, true)
	if err := f.Write(changes()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "+ [db.host]") {
		t.Errorf("expected added key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "- [db.port]") {
		t.Errorf("expected removed key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "~ [app.env]") {
		t.Errorf("expected modified key in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Summary:") {
		t.Errorf("expected summary line in output, got:\n%s", out)
	}
}

func TestWriteText_Summary(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatText, true)
	_ = f.Write(changes())
	out := buf.String()
	if !strings.Contains(out, "1 added, 1 removed, 1 modified") {
		t.Errorf("unexpected summary: %s", out)
	}
}

func TestWriteJSON_Output(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON, true)
	if err := f.Write(changes()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.HasPrefix(out, "[") || !strings.Contains(out, "]") {
		t.Errorf("expected JSON array, got:\n%s", out)
	}
	if !strings.Contains(out, `"db.host"`) {
		t.Errorf("expected key in JSON output, got:\n%s", out)
	}
	if !strings.Contains(out, `"added"`) {
		t.Errorf("expected type 'added' in JSON output, got:\n%s", out)
	}
}

func TestWriteJSON_EmptyChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON, true)
	if err := f.Write([]diff.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "[]" {
		t.Errorf("expected empty JSON array, got: %s", out)
	}
}
