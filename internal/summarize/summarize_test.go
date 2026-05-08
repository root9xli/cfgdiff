package summarize_test

import (
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/diff"
	"github.com/user/cfgdiff/internal/summarize"
)

var sampleChanges = []diff.Change{
	{Key: "host", Type: diff.Added, NewValue: "localhost"},
	{Key: "port", Type: diff.Modified, OldValue: "8080", NewValue: "9090"},
	{Key: "debug", Type: diff.Removed, OldValue: "true"},
	{Key: "port", Type: diff.Modified, OldValue: "9090", NewValue: "9091"},
}

func TestCompute_Counts(t *testing.T) {
	s := summarize.Compute(sampleChanges, summarize.DefaultOptions())
	if s.Total != 4 {
		t.Errorf("expected Total=4, got %d", s.Total)
	}
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected Removed=1, got %d", s.Removed)
	}
	if s.Modified != 2 {
		t.Errorf("expected Modified=2, got %d", s.Modified)
	}
}

func TestCompute_TopKeys(t *testing.T) {
	s := summarize.Compute(sampleChanges, summarize.Options{TopN: 2})
	if len(s.TopKeys) == 0 {
		t.Fatal("expected at least one top key")
	}
	// "port" appears twice so must be first
	if s.TopKeys[0] != "port" {
		t.Errorf("expected first top key to be 'port', got %q", s.TopKeys[0])
	}
}

func TestCompute_Empty(t *testing.T) {
	s := summarize.Compute(nil, summarize.DefaultOptions())
	if s.Total != 0 {
		t.Errorf("expected Total=0 for empty input, got %d", s.Total)
	}
	if len(s.TopKeys) != 0 {
		t.Errorf("expected no top keys for empty input")
	}
}

func TestFormat_String(t *testing.T) {
	s := summarize.Stats{Total: 3, Added: 1, Removed: 1, Modified: 1}
	out := summarize.Format(s)
	if !strings.Contains(out, "total=3") {
		t.Errorf("Format output missing total: %q", out)
	}
	if !strings.Contains(out, "added=1") {
		t.Errorf("Format output missing added: %q", out)
	}
}

func TestCompute_ZeroTopN_UsesDefault(t *testing.T) {
	s := summarize.Compute(sampleChanges, summarize.Options{TopN: 0})
	// should not panic and should return some keys
	if s.Total == 0 {
		t.Error("expected non-zero total")
	}
}
