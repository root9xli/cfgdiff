package diff

import (
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": "8080"}
	target := map[string]interface{}{"host": "localhost", "port": "8080"}

	result := Compare(base, target)

	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
	if result.Summary() != "No differences found." {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestCompare_Added(t *testing.T) {
	base := map[string]interface{}{"host": "localhost"}
	target := map[string]interface{}{"host": "localhost", "port": "9090"}

	result := Compare(base, target)

	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].Key != "port" {
		t.Errorf("expected key 'port', got '%s'", result.Changes[0].Key)
	}
}

func TestCompare_Removed(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "debug": "true"}
	target := map[string]interface{}{"host": "localhost"}

	result := Compare(base, target)

	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", result.Changes[0].Type)
	}
}

func TestCompare_Modified(t *testing.T) {
	base := map[string]interface{}{"host": "localhost", "port": "8080"}
	target := map[string]interface{}{"host": "localhost", "port": "9090"}

	result := Compare(base, target)

	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != Modified {
		t.Errorf("expected Modified, got %s", c.Type)
	}
	if c.OldValue != "8080" || c.NewValue != "9090" {
		t.Errorf("unexpected values: old=%v new=%v", c.OldValue, c.NewValue)
	}
}

func TestCompare_SummaryFormat(t *testing.T) {
	base := map[string]interface{}{"a": "1", "b": "2", "c": "3"}
	target := map[string]interface{}{"a": "1", "b": "99", "d": "4"}

	result := Compare(base, target)

	expected := "1 added, 1 removed, 1 modified"
	if result.Summary() != expected {
		t.Errorf("expected summary %q, got %q", expected, result.Summary())
	}
}

func TestCompare_SortedKeys(t *testing.T) {
	base := map[string]interface{}{}
	target := map[string]interface{}{"zebra": "1", "apple": "2", "mango": "3"}

	result := Compare(base, target)

	if result.Changes[0].Key != "apple" || result.Changes[1].Key != "mango" || result.Changes[2].Key != "zebra" {
		t.Errorf("changes not sorted by key: %v", result.Changes)
	}
}
