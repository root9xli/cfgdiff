package classify_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/cfgdiff/internal/classify"
	"github.com/user/cfgdiff/internal/diff"
)

func TestFprintResults_NoResults(t *testing.T) {
	var buf bytes.Buffer
	classify.FprintResults(&buf, []classify.Result{})
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes' message, got: %s", buf.String())
	}
}

func TestFprintResults_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	c := classify.New(classify.DefaultOptions())
	results := c.Apply([]diff.Change{
		{Key: "db.password", Type: "modified", OldValue: "a", NewValue: "b"},
	})
	classify.FprintResults(&buf, results)
	out := buf.String()
	for _, col := range []string{"KEY", "TYPE", "DOMAIN", "SEVERITY"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestFprintResults_ContainsRowData(t *testing.T) {
	var buf bytes.Buffer
	c := classify.New(classify.DefaultOptions())
	results := c.Apply([]diff.Change{
		{Key: "server.host", Type: "added", OldValue: "", NewValue: "prod"},
	})
	classify.FprintResults(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "server.host") {
		t.Errorf("expected key in output")
	}
	if !strings.Contains(out, "high") {
		t.Errorf("expected severity 'high' in output")
	}
	if !strings.Contains(out, "server") {
		t.Errorf("expected domain 'server' in output")
	}
}
