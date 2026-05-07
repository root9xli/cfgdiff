package template_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/diff"
	"github.com/user/cfgdiff/internal/template"
)

var sampleChanges = []diff.Change{
	{Key: "host", Type: diff.Added, NewValue: "localhost"},
	{Key: "port", Type: diff.Removed, OldValue: "8080"},
	{Key: "timeout", Type: diff.Modified, OldValue: "30", NewValue: "60"},
}

func TestNew_ValidTemplate(t *testing.T) {
	_, err := template.New("{{ range .Changes }}{{ .Key }}\n{{ end }}")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := template.New("{{ range .Changes }")
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}

func TestRender_ListsKeys(t *testing.T) {
	r, err := template.New("{{ range .Changes }}{{ .Key }}\n{{ end }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, sampleChanges); err != nil {
		t.Fatalf("render error: %v", err)
	}

	out := buf.String()
	for _, key := range []string{"host", "port", "timeout"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain %q, got:\n%s", key, out)
		}
	}
}

func TestRender_SummaryCounts(t *testing.T) {
	src := "added={{ .Summary.Added }} removed={{ .Summary.Removed }} modified={{ .Summary.Modified }}"
	r, err := template.New(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, sampleChanges); err != nil {
		t.Fatalf("render error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "added=1") {
		t.Errorf("expected added=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "removed=1") {
		t.Errorf("expected removed=1 in output, got: %s", out)
	}
	if !strings.Contains(out, "modified=1") {
		t.Errorf("expected modified=1 in output, got: %s", out)
	}
}

func TestRender_EmptyChanges(t *testing.T) {
	r, err := template.New("total={{ .Summary.Total }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, []diff.Change{}); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if got := buf.String(); got != "total=0" {
		t.Errorf("expected 'total=0', got %q", got)
	}
}

func TestRender_RenderedTimestampPresent(t *testing.T) {
	r, err := template.New("{{ .Rendered.Year }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	if err := r.Render(&buf, nil); err != nil {
		t.Fatalf("render error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty rendered timestamp")
	}
}
