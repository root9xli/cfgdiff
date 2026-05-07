package filter_test

import (
	"testing"

	"github.com/yourorg/cfgdiff/internal/diff"
	"github.com/yourorg/cfgdiff/internal/filter"
)

var sampleChanges = []diff.Change{
	{Key: "db.host", Type: "modified", OldValue: "localhost", NewValue: "prod-db"},
	{Key: "db.port", Type: "added", OldValue: "", NewValue: "5432"},
	{Key: "app.debug", Type: "removed", OldValue: "true", NewValue: ""},
	{Key: "app.name", Type: "modified", OldValue: "myapp", NewValue: "myapp-v2"},
	{Key: "cache.ttl", Type: "added", OldValue: "", NewValue: "300"},
}

func TestApply_NoFilter_ReturnsAll(t *testing.T) {
	f := filter.New(filter.Options{})
	result := f.Apply(sampleChanges)
	if len(result) != len(sampleChanges) {
		t.Errorf("expected %d changes, got %d", len(sampleChanges), len(result))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	f := filter.New(filter.Options{Prefix: "db."})
	result := f.Apply(sampleChanges)
	if len(result) != 2 {
		t.Errorf("expected 2 changes, got %d", len(result))
	}
	for _, c := range result {
		if c.Key[:3] != "db." {
			t.Errorf("unexpected key %q with prefix filter", c.Key)
		}
	}
}

func TestApply_ContainsFilter(t *testing.T) {
	f := filter.New(filter.Options{Contains: "app"})
	result := f.Apply(sampleChanges)
	if len(result) != 2 {
		t.Errorf("expected 2 changes, got %d", len(result))
	}
}

func TestApply_TypeFilter_Added(t *testing.T) {
	f := filter.New(filter.Options{Types: []string{"added"}})
	result := f.Apply(sampleChanges)
	if len(result) != 2 {
		t.Errorf("expected 2 added changes, got %d", len(result))
	}
	for _, c := range result {
		if c.Type != "added" {
			t.Errorf("expected type 'added', got %q", c.Type)
		}
	}
}

func TestApply_TypeFilter_CaseInsensitive(t *testing.T) {
	f := filter.New(filter.Options{Types: []string{"MODIFIED"}})
	result := f.Apply(sampleChanges)
	if len(result) != 2 {
		t.Errorf("expected 2 modified changes, got %d", len(result))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	f := filter.New(filter.Options{
		Prefix: "db.",
		Types:  []string{"added"},
	})
	result := f.Apply(sampleChanges)
	if len(result) != 1 {
		t.Errorf("expected 1 change, got %d", len(result))
	}
	if result[0].Key != "db.port" {
		t.Errorf("expected key 'db.port', got %q", result[0].Key)
	}
}

func TestApply_NoMatches(t *testing.T) {
	f := filter.New(filter.Options{Prefix: "nonexistent."})
	result := f.Apply(sampleChanges)
	if result != nil && len(result) != 0 {
		t.Errorf("expected no results, got %d", len(result))
	}
}
