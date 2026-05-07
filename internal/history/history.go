package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/user/cfgdiff/internal/diff"
)

// Entry represents a single recorded diff history entry.
type Entry struct {
	ID        string          `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	FileA     string          `json:"file_a"`
	FileB     string          `json:"file_b"`
	Changes   []diff.Change   `json:"changes"`
	Summary   Summary         `json:"summary"`
}

// Summary holds aggregated counts for an entry.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
}

// Store manages persisted diff history.
type Store struct {
	dir string
}

// NewStore creates a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Record saves a new history entry derived from the given changes.
func (s *Store) Record(fileA, fileB string, changes []diff.Change) (*Entry, error) {
	sum := buildSummary(changes)
	entry := &Entry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now().UTC(),
		FileA:     fileA,
		FileB:     fileB,
		Changes:   changes,
		Summary:   sum,
	}
	path := filepath.Join(s.dir, entry.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(entry); err != nil {
		return nil, fmt.Errorf("history: encode: %w", err)
	}
	return entry, nil
}

// List returns all entries sorted by timestamp descending.
func (s *Store) List() ([]Entry, error) {
	glob := filepath.Join(s.dir, "*.json")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, m := range matches {
		e, err := loadEntry(m)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	return entries, nil
}

// Get retrieves a single entry by ID.
func (s *Store) Get(id string) (*Entry, error) {
	return loadEntry(filepath.Join(s.dir, id+".json"))
}

func loadEntry(path string) (*Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("history: entry not found")
		}
		return nil, fmt.Errorf("history: open: %w", err)
	}
	defer f.Close()
	var e Entry
	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, fmt.Errorf("history: decode: %w", err)
	}
	return &e, nil
}

func buildSummary(changes []diff.Change) Summary {
	var s Summary
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		}
	}
	return s
}
