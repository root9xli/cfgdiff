package history

import (
	"os"
	"testing"

	"github.com/user/cfgdiff/internal/diff"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir, err := os.MkdirTemp("", "cfgdiff-history-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

var sampleChanges = []diff.Change{
	{Type: diff.Added, Key: "NEW_KEY", NewValue: "v"},
	{Type: diff.Removed, Key: "OLD_KEY", OldValue: "v"},
	{Type: diff.Modified, Key: "MOD_KEY", OldValue: "a", NewValue: "b"},
}

func TestRecord_CreatesEntry(t *testing.T) {
	s := newTestStore(t)
	e, err := s.Record("a.env", "b.env", sampleChanges)
	if err != nil {
		t.Fatalf("Record: %v", err)
	}
	if e.ID == "" {
		t.Error("expected non-empty ID")
	}
	if e.Summary.Added != 1 || e.Summary.Removed != 1 || e.Summary.Modified != 1 {
		t.Errorf("unexpected summary: %+v", e.Summary)
	}
}

func TestList_Empty(t *testing.T) {
	s := newTestStore(t)
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestList_MultipleEntries(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 3; i++ {
		if _, err := s.Record("a.env", "b.env", sampleChanges); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestGet_Found(t *testing.T) {
	s := newTestStore(t)
	recorded, _ := s.Record("a.env", "b.env", sampleChanges)
	fetched, err := s.Get(recorded.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if fetched.ID != recorded.ID {
		t.Errorf("ID mismatch: got %s want %s", fetched.ID, recorded.ID)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing entry")
	}
}

func TestBuildSummary(t *testing.T) {
	sum := buildSummary(sampleChanges)
	if sum.Added != 1 || sum.Removed != 1 || sum.Modified != 1 {
		t.Errorf("unexpected summary: %+v", sum)
	}
}
