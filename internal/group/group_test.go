package group

import (
	"testing"

	"github.com/cfgdiff/internal/diff"
)

var sampleChanges = []diff.Change{
	{Key: "database.host", Type: diff.Modified, OldValue: "localhost", NewValue: "db.prod"},
	{Key: "database.port", Type: diff.Added, NewValue: "5432"},
	{Key: "server.port", Type: diff.Modified, OldValue: "8080", NewValue: "443"},
	{Key: "debug", Type: diff.Removed, OldValue: "true"},
}

func TestApply_GroupsByPrefix(t *testing.T) {
	results := Apply(sampleChanges, DefaultOptions())
	if len(results) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(results))
	}
	// Results are sorted; expect (root), database, server
	if results[0].Prefix != "(root)" {
		t.Errorf("expected first prefix '(root)', got %q", results[0].Prefix)
	}
	if results[1].Prefix != "database" {
		t.Errorf("expected second prefix 'database', got %q", results[1].Prefix)
	}
	if results[2].Prefix != "server" {
		t.Errorf("expected third prefix 'server', got %q", results[2].Prefix)
	}
}

func TestApply_RootGroup_HasOneChange(t *testing.T) {
	results := Apply(sampleChanges, DefaultOptions())
	for _, r := range results {
		if r.Prefix == "(root)" && len(r.Changes) != 1 {
			t.Errorf("expected 1 root change, got %d", len(r.Changes))
		}
	}
}

func TestApply_DatabaseGroup_HasTwoChanges(t *testing.T) {
	results := Apply(sampleChanges, DefaultOptions())
	for _, r := range results {
		if r.Prefix == "database" && len(r.Changes) != 2 {
			t.Errorf("expected 2 database changes, got %d", len(r.Changes))
		}
	}
}

func TestApply_EmptyChanges(t *testing.T) {
	results := Apply([]diff.Change{}, DefaultOptions())
	if len(results) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(results))
	}
}

func TestApply_CustomDepth(t *testing.T) {
	changes := []diff.Change{
		{Key: "a.b.c", Type: diff.Added, NewValue: "1"},
		{Key: "a.b.d", Type: diff.Added, NewValue: "2"},
		{Key: "a.x.y", Type: diff.Added, NewValue: "3"},
	}
	opts := Options{Separator: ".", Depth: 2}
	results := Apply(changes, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 groups at depth 2, got %d", len(results))
	}
	if results[0].Prefix != "a.b" {
		t.Errorf("expected prefix 'a.b', got %q", results[0].Prefix)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	changes := []diff.Change{
		{Key: "db/host", Type: diff.Modified, OldValue: "a", NewValue: "b"},
		{Key: "db/port", Type: diff.Added, NewValue: "5432"},
	}
	opts := Options{Separator: "/", Depth: 1}
	results := Apply(changes, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 group, got %d", len(results))
	}
	if results[0].Prefix != "db" {
		t.Errorf("expected prefix 'db', got %q", results[0].Prefix)
	}
}
