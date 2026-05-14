package annotate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/annotate"
	"github.com/user/cfgdiff/internal/diff"
)

var sampleChanges = []diff.Change{
	{Key: "APP_HOST", Type: diff.Added, NewValue: "localhost"},
	{Key: "APP_PORT", Type: diff.Removed, OldValue: "8080"},
	{Key: "DB_PASS", Type: diff.Modified, OldValue: "old", NewValue: "new"},
}

func TestApply_ReturnsOneAnnotationPerChange(t *testing.T) {
	a := annotate.New(annotate.DefaultOptions())
	result := a.Apply(sampleChanges)
	if len(result) != len(sampleChanges) {
		t.Fatalf("expected %d annotations, got %d", len(sampleChanges), len(result))
	}
}

func TestApply_AddedNote(t *testing.T) {
	a := annotate.New(annotate.DefaultOptions())
	result := a.Apply([]diff.Change{
		{Key: "FOO", Type: diff.Added, NewValue: "bar"},
	})
	if !strings.Contains(result[0].Note, "introduced") {
		t.Errorf("expected 'introduced' in note, got: %s", result[0].Note)
	}
}

func TestApply_RemovedNote(t *testing.T) {
	a := annotate.New(annotate.DefaultOptions())
	result := a.Apply([]diff.Change{
		{Key: "FOO", Type: diff.Removed, OldValue: "bar"},
	})
	if !strings.Contains(result[0].Note, "removed") {
		t.Errorf("expected 'removed' in note, got: %s", result[0].Note)
	}
}

func TestApply_ModifiedNote(t *testing.T) {
	a := annotate.New(annotate.DefaultOptions())
	result := a.Apply([]diff.Change{
		{Key: "FOO", Type: diff.Modified, OldValue: "a", NewValue: "b"},
	})
	if !strings.Contains(result[0].Note, "changed") {
		t.Errorf("expected 'changed' in note, got: %s", result[0].Note)
	}
}

func TestApply_CustomNoteOverridesDefault(t *testing.T) {
	opts := annotate.Options{
		CustomNotes: map[string]string{"DB_": "database credential change"},
	}
	a := annotate.New(opts)
	result := a.Apply([]diff.Change{
		{Key: "DB_PASS", Type: diff.Modified, OldValue: "x", NewValue: "y"},
	})
	if result[0].Note != "database credential change" {
		t.Errorf("unexpected note: %s", result[0].Note)
	}
}

func TestFprint_NoAnnotations(t *testing.T) {
	var buf bytes.Buffer
	annotate.Fprint(&buf, nil)
	if !strings.Contains(buf.String(), "no annotations") {
		t.Errorf("expected 'no annotations', got: %s", buf.String())
	}
}

func TestFprint_ContainsHeader(t *testing.T) {
	a := annotate.New(annotate.DefaultOptions())
	anns := a.Apply(sampleChanges)
	var buf bytes.Buffer
	annotate.Fprint(&buf, anns)
	out := buf.String()
	for _, header := range []string{"TYPE", "KEY", "NOTE"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}
