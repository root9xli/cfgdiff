package truncate_test

import (
	"strings"
	"testing"

	"github.com/cfgdiff/cfgdiff/internal/truncate"
)

func TestValue_ShortString_Unchanged(t *testing.T) {
	tr := truncate.New(truncate.DefaultOptions())
	input := "short"
	if got := tr.Value(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestValue_LongString_Truncated(t *testing.T) {
	opts := truncate.Options{MaxLen: 10, Suffix: "...", Enabled: true}
	tr := truncate.New(opts)
	input := "this is a very long configuration value"
	got := tr.Value(input)
	if len(got) > 10 {
		t.Errorf("expected len <= 10, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestValue_Disabled_NoTruncation(t *testing.T) {
	opts := truncate.Options{MaxLen: 5, Suffix: "...", Enabled: false}
	tr := truncate.New(opts)
	input := "this should not be truncated at all"
	if got := tr.Value(input); got != input {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestApply_TruncatesMapValues(t *testing.T) {
	opts := truncate.Options{MaxLen: 8, Suffix: "...", Enabled: true}
	tr := truncate.New(opts)
	m := map[string]string{
		"key1": "short",
		"key2": "a very long value here",
	}
	out := tr.Apply(m)
	if out["key1"] != "short" {
		t.Errorf("key1 should be unchanged, got %q", out["key1"])
	}
	if len(out["key2"]) > 8 {
		t.Errorf("key2 should be truncated, got %q", out["key2"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	opts := truncate.Options{MaxLen: 5, Suffix: "...", Enabled: true}
	tr := truncate.New(opts)
	m := map[string]string{"k": "original long value"}
	_ = tr.Apply(m)
	if m["k"] != "original long value" {
		t.Error("original map was mutated")
	}
}

func TestSummary_NoTruncation(t *testing.T) {
	tr := truncate.New(truncate.DefaultOptions())
	m := map[string]string{"a": "short"}
	out := tr.Apply(m)
	if got := tr.Summary(m, out); got != "no values truncated" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithTruncation(t *testing.T) {
	opts := truncate.Options{MaxLen: 5, Suffix: "...", Enabled: true}
	tr := truncate.New(opts)
	m := map[string]string{"a": "longvalue", "b": "hi"}
	out := tr.Apply(m)
	summary := tr.Summary(m, out)
	if !strings.Contains(summary, "1 value(s) truncated") {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestNew_ZeroMaxLen_UsesDefault(t *testing.T) {
	opts := truncate.Options{MaxLen: 0, Suffix: "...", Enabled: true}
	tr := truncate.New(opts)
	long := strings.Repeat("x", 200)
	got := tr.Value(long)
	if len(got) > truncate.DefaultMaxLen {
		t.Errorf("expected default max len %d, got %d", truncate.DefaultMaxLen, len(got))
	}
}
