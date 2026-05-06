package watch

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// FileState holds the last known state of a watched file.
type FileState struct {
	Path    string
	Checksum string
	ModTime time.Time
}

// Watcher polls files for changes at a given interval.
type Watcher struct {
	files    []string
	interval time.Duration
	states   map[string]FileState
	OnChange func(path string, prev, curr FileState)
}

// New creates a Watcher for the given file paths and poll interval.
func New(files []string, interval time.Duration) *Watcher {
	return &Watcher{
		files:    files,
		interval: interval,
		states:   make(map[string]FileState),
	}
}

// Start begins polling. It blocks until the done channel is closed.
func (w *Watcher) Start(done <-chan struct{}) error {
	if err := w.snapshot(); err != nil {
		return err
	}
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return nil
		case <-ticker.C:
			w.poll()
		}
	}
}

func (w *Watcher) snapshot() error {
	for _, path := range w.files {
		state, err := fileState(path)
		if err != nil {
			return fmt.Errorf("watch: initial snapshot failed for %s: %w", path, err)
		}
		w.states[path] = state
	}
	return nil
}

func (w *Watcher) poll() {
	for _, path := range w.files {
		curr, err := fileState(path)
		if err != nil {
			continue
		}
		prev, known := w.states[path]
		if !known || curr.Checksum != prev.Checksum {
			w.states[path] = curr
			if w.OnChange != nil {
				w.OnChange(path, prev, curr)
			}
		}
	}
}

func fileState(path string) (FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileState{}, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return FileState{}, err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return FileState{}, err
	}
	return FileState{
		Path:    path,
		Checksum: fmt.Sprintf("%x", h.Sum(nil)),
		ModTime: info.ModTime(),
	}, nil
}
