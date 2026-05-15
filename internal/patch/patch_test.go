package patch

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cfgdiff/internal/diff"
)

var baseConfig = map[string]string{
	"HOST":  "localhost",
	"PORT":  "5432",
	"DEBUG": "false",
}

func TestApplyForward_Added(t *testing.T) {
	changes := []diff.Change{{Type: diff.Added, Key: "TIMEOUT", NewValue: "30s"}}
	p := New(changes, Forward)
	out, err := p.Apply(baseConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TIMEOUT"] != "30s" {
		t.Errorf("expected TIMEOUT=30s, got %q", out["TIMEOUT"])
	}
}

func TestApplyForward_Removed(t *testing.T) {
	changes := []diff.Change{{Type: diff.Removed, Key: "DEBUG", OldValue: "false"}}
	p := New(changes, Forward)
	out, err := p.Apply(baseConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DEBUG"]; ok {
		t.Error("expected DEBUG to be removed")
	}
}

func TestApplyForward_Modified(t *testing.T) {
	changes := []diff.Change{{Type: diff.Modified, Key: "PORT", OldValue: "5432", NewValue: "5433"}}
	p := New(changes, Forward)
	out, err := p.Apply(baseConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "5433" {
		t.Errorf("expected PORT=5433, got %q", out["PORT"])
	}
}

func TestApplyReverse_Modified(t *testing.T) {
	modified := map[string]string{"HOST": "localhost", "PORT": "5433", "DEBUG": "false"}
	changes := []diff.Change{{Type: diff.Modified, Key: "PORT", OldValue: "5432", NewValue: "5433"}}
	p := New(changes, Reverse)
	out, err := p.Apply(modified)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "5432" {
		t.Errorf("expected PORT reverted to 5432, got %q", out["PORT"])
	}
}

func TestApplyForward_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"HOST": "localhost"}
	changes := []diff.Change{{Type: diff.Added, Key: "PORT", NewValue: "8080"}}
	p := New(changes, Forward)
	_, _ = p.Apply(original)
	if _, ok := original["PORT"]; ok {
		t.Error("Apply mutated the original config")
	}
}

func TestApplyForward_MissingKeyForModify(t *testing.T) {
	changes := []diff.Change{{Type: diff.Modified, Key: "MISSING", OldValue: "x", NewValue: "y"}}
	p := New(changes, Forward)
	_, err := p.Apply(baseConfig)
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

// TestApplyReverse_Added verifies that reversing an "added" change removes the key.
func TestApplyReverse_Added(t *testing.T) {
	withExtra := map[string]string{"HOST": "localhost", "PORT": "5432", "DEBUG": "false", "TIMEOUT": "30s"}
	changes := []diff.Change{{Type: diff.Added, Key: "TIMEOUT", NewValue: "30s"}}
	p := New(changes, Reverse)
	out, err := p.Apply(withExtra)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["TIMEOUT"]; ok {
		t.Error("expected TIMEOUT to be removed when reversing an Added change")
	}
}

func TestFprintPatch_Output(t *testing.T) {
	changes := []diff.Change{
		{Type: diff.Added, Key: "NEW_KEY", NewValue: "val"},
		{Type: diff.Removed, Key: "OLD_KEY", OldValue: "old"},
		{Type: diff.Modified, Key: "HOST", OldValue: "localhost", NewValue: "prod.host"},
	}
	p := New(changes, Forward)
	var buf bytes.Buffer
	FprintPatch(&buf, p)
	out := buf.String()
	if !strings.Contains(out, "forward") {
		t.Error("expected direction 'forward' in output")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("expected NEW_KEY in output")
	}
	if !strings.Contains(out, "->") {
		t.Error("expected '->' separator for modified key")
	}
}
