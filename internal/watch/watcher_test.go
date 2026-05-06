package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestFileState_Checksum(t *testing.T) {
	path := writeTempFile(t, "key: value\n")
	state, err := fileState(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if state.Path != path {
		t.Errorf("expected path %s, got %s", path, state.Path)
	}
}

func TestFileState_ChangesOnWrite(t *testing.T) {
	path := writeTempFile(t, "key: original\n")
	s1, _ := fileState(path)
	if err := os.WriteFile(path, []byte("key: modified\n"), 0644); err != nil {
		t.Fatal(err)
	}
	s2, _ := fileState(path)
	if s1.Checksum == s2.Checksum {
		t.Error("expected checksums to differ after file change")
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempFile(t, "a: 1\n")
	w := New([]string{path}, 20*time.Millisecond)

	changed := make(chan string, 1)
	w.OnChange = func(p string, prev, curr FileState) {
		changed <- p
	}

	done := make(chan struct{})
	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(done)
	}()

	time.Sleep(40 * time.Millisecond)
	if err := os.WriteFile(path, []byte("a: 2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	select {
	case p := <-changed:
		if p != path {
			t.Errorf("expected %s, got %s", path, p)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for change event")
	}
	close(done)
	if err := <-errCh; err != nil {
		t.Errorf("watcher error: %v", err)
	}
}

func TestWatcher_NoSpuriousEvents(t *testing.T) {
	path := writeTempFile(t, "stable: true\n")
	w := New([]string{path}, 20*time.Millisecond)

	count := 0
	w.OnChange = func(_ string, _, _ FileState) { count++ }

	done := make(chan struct{})
	go func() { w.Start(done) }()
	time.Sleep(120 * time.Millisecond)
	close(done)

	if count != 0 {
		t.Errorf("expected 0 change events for stable file, got %d", count)
	}
}
