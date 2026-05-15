package rename

import (
	"bytes"
	"testing"

	"github.com/cfgdiff/internal/diff"
)

func TestApplyToMap_RenamesKey(t *testing.T) {
	rules := []Rule{{From: "db_host", To: "database_host"}}
	r := New(rules, DefaultOptions())
	in := map[string]string{"db_host": "localhost", "port": "5432"}
	out := r.ApplyToMap(in)
	if _, ok := out["database_host"]; !ok {
		t.Error("expected key 'database_host' in output")
	}
	if _, ok := out["db_host"]; ok {
		t.Error("old key 'db_host' should not appear in output")
	}
	if out["port"] != "5432" {
		t.Error("unrelated key should be preserved")
	}
}

func TestApplyToMap_DoesNotMutateOriginal(t *testing.T) {
	rules := []Rule{{From: "old", To: "new"}}
	r := New(rules, DefaultOptions())
	in := map[string]string{"old": "val"}
	_ = r.ApplyToMap(in)
	if _, ok := in["old"]; !ok {
		t.Error("original map should not be mutated")
	}
}

func TestApplyToChanges_RenamesKey(t *testing.T) {
	rules := []Rule{{From: "app_port", To: "server_port"}}
	r := New(rules, DefaultOptions())
	changes := []diff.Change{
		{Key: "app_port", Type: diff.Modified, OldValue: "8080", NewValue: "9090"},
		{Key: "log_level", Type: diff.Added, NewValue: "debug"},
	}
	out := r.ApplyToChanges(changes)
	if out[0].Key != "server_port" {
		t.Errorf("expected 'server_port', got %q", out[0].Key)
	}
	if out[1].Key != "log_level" {
		t.Error("unrelated change key should be unchanged")
	}
}

func TestApplyToChanges_CaseInsensitive(t *testing.T) {
	rules := []Rule{{From: "DB_HOST", To: "database_host"}}
	opts := Options{CaseSensitive: false}
	r := New(rules, opts)
	changes := []diff.Change{{Key: "db_host", Type: diff.Added, NewValue: "localhost"}}
	out := r.ApplyToChanges(changes)
	if out[0].Key != "database_host" {
		t.Errorf("expected 'database_host', got %q", out[0].Key)
	}
}

func TestLoadRules_ValidPairs(t *testing.T) {
	pairs := []string{"old_key=new_key", "foo=bar"}
	rules, err := LoadRules(pairs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "old_key" || rules[0].To != "new_key" {
		t.Errorf("unexpected rule: %+v", rules[0])
	}
}

func TestLoadRules_SkipsInvalid(t *testing.T) {
	pairs := []string{"=newval", "oldval=", "valid=pair"}
	rules, _ := LoadRules(pairs)
	if len(rules) != 1 {
		t.Errorf("expected 1 valid rule, got %d", len(rules))
	}
}

func TestFprintRules_NoRules(t *testing.T) {
	var buf bytes.Buffer
	FprintRules(&buf, []Rule{})
	if buf.String() == "" {
		t.Error("expected non-empty output for empty rules")
	}
}

func TestFprintRules_WithRules(t *testing.T) {
	var buf bytes.Buffer
	rules := []Rule{{From: "a", To: "b"}, {From: "x", To: "y"}}
	FprintRules(&buf, rules)
	out := buf.String()
	if !containsStr(out, "a") || !containsStr(out, "b") {
		t.Error("expected rule entries in output")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
