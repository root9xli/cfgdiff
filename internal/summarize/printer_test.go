package summarize_test

import (
	"strings"
	"testing"

	"bytes"

	"github.com/user/cfgdiff/internal/summarize"
)

func TestPrint_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.Stats{Total: 2, Added: 1, Removed: 1}
	summarize.Print(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "Change Summary") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestPrint_ShowsCounts(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.Stats{Total: 5, Added: 2, Removed: 1, Modified: 2}
	summarize.Print(&buf, s)
	out := buf.String()
	for _, want := range []string{"5", "2", "1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}

func TestPrint_ShowsTopKeys(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.Stats{
		Total:   1,
		Added:   1,
		TopKeys: []string{"database_url", "secret_key"},
	}
	summarize.Print(&buf, s)
	out := buf.String()
	if !strings.Contains(out, "database_url") {
		t.Errorf("expected top key in output:\n%s", out)
	}
}

func TestPrint_NoTopKeys_OmitsLine(t *testing.T) {
	var buf bytes.Buffer
	s := summarize.Stats{Total: 0}
	summarize.Print(&buf, s)
	out := buf.String()
	if strings.Contains(out, "Top keys") {
		t.Errorf("expected no 'Top keys' line when empty, got:\n%s", out)
	}
}
