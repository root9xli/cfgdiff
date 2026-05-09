package classify_test

import (
	"testing"

	"github.com/user/cfgdiff/internal/classify"
	"github.com/user/cfgdiff/internal/diff"
)

func sampleChanges() []diff.Change {
	return []diff.Change{
		{Key: "db.password", Type: "modified", OldValue: "old", NewValue: "new"},
		{Key: "server.host", Type: "modified", OldValue: "localhost", NewValue: "prod.host"},
		{Key: "app.name", Type: "added", OldValue: "", NewValue: "myapp"},
		{Key: "auth.token", Type: "removed", OldValue: "abc", NewValue: ""},
	}
}

func TestApply_ReturnsSameCount(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply(sampleChanges())
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}

func TestApply_CriticalSeverity(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply(sampleChanges())
	// db.password should be critical
	if results[0].Severity != classify.SeverityCritical {
		t.Errorf("expected critical, got %s", results[0].Severity)
	}
	// auth.token should be critical
	if results[3].Severity != classify.SeverityCritical {
		t.Errorf("expected critical, got %s", results[3].Severity)
	}
}

func TestApply_HighSeverity(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply(sampleChanges())
	if results[1].Severity != classify.SeverityHigh {
		t.Errorf("expected high, got %s", results[1].Severity)
	}
}

func TestApply_LowSeverity(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply(sampleChanges())
	if results[2].Severity != classify.SeverityLow {
		t.Errorf("expected low, got %s", results[2].Severity)
	}
}

func TestApply_DomainFromDotNotation(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply(sampleChanges())
	if results[0].Domain != "db" {
		t.Errorf("expected domain 'db', got %s", results[0].Domain)
	}
}

func TestApply_DomainFromUnderscore(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	changes := []diff.Change{
		{Key: "APP_NAME", Type: "added", OldValue: "", NewValue: "x"},
	}
	results := c.Apply(changes)
	if results[0].Domain != "app" {
		t.Errorf("expected domain 'app', got %s", results[0].Domain)
	}
}

func TestApply_EmptyChanges(t *testing.T) {
	c := classify.New(classify.DefaultOptions())
	results := c.Apply([]diff.Change{})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	opts := classify.Options{
		CriticalPatterns: []string{"private"},
		HighPatterns:     []string{},
	}
	c := classify.New(opts)
	changes := []diff.Change{
		{Key: "private_key", Type: "modified", OldValue: "a", NewValue: "b"},
	}
	results := c.Apply(changes)
	if results[0].Severity != classify.SeverityCritical {
		t.Errorf("expected critical, got %s", results[0].Severity)
	}
}
