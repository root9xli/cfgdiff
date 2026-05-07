package history

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/cfgdiff/internal/diff"
)

func sampleEntry() *Entry {
	return &Entry{
		ID:        "123456",
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		FileA:     "prod.env",
		FileB:     "staging.env",
		Changes:   sampleChanges,
		Summary:   Summary{Added: 1, Removed: 1, Modified: 1},
	}
}

func TestPrintList_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	PrintList(&buf, nil)
	if !strings.Contains(buf.String(), "No history") {
		t.Errorf("expected no-entries message, got: %s", buf.String())
	}
}

func TestPrintList_WithEntries(t *testing.T) {
	var buf bytes.Buffer
	PrintList(&buf, []Entry{*sampleEntry()})
	out := buf.String()
	for _, want := range []string{"123456", "prod.env", "staging.env", "1", "2024-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestPrintEntry_Detail(t *testing.T) {
	var buf bytes.Buffer
	PrintEntry(&buf, sampleEntry())
	out := buf.String()
	for _, want := range []string{"123456", "prod.env", "staging.env", "NEW_KEY", "OLD_KEY", "MOD_KEY"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestPrintEntry_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	e := &Entry{ID: "0", FileA: "a", FileB: "b", Changes: []diff.Change{}}
	PrintEntry(&buf, e)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}
