package watch

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrintChange_FirstSeen(t *testing.T) {
	var buf bytes.Buffer
	curr := FileState{
		Path:     "/etc/app.yaml",
		Checksum: "abc123",
		ModTime:  time.Now(),
	}
	PrintChange(&buf, "/etc/app.yaml", FileState{}, curr)
	out := buf.String()
	if !strings.Contains(out, "DETECTED") {
		t.Errorf("expected DETECTED in output, got: %s", out)
	}
	if !strings.Contains(out, "abc123") {
		t.Errorf("expected checksum in output, got: %s", out)
	}
}

func TestPrintChange_Modified(t *testing.T) {
	var buf bytes.Buffer
	prev := FileState{Path: "/etc/app.yaml", Checksum: "aaa", ModTime: time.Now()}
	curr := FileState{Path: "/etc/app.yaml", Checksum: "bbb", ModTime: time.Now()}
	PrintChange(&buf, "/etc/app.yaml", prev, curr)
	out := buf.String()
	if !strings.Contains(out, "CHANGED") {
		t.Errorf("expected CHANGED in output, got: %s", out)
	}
	if !strings.Contains(out, "aaa") || !strings.Contains(out, "bbb") {
		t.Errorf("expected both checksums in output, got: %s", out)
	}
}
