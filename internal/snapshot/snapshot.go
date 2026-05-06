package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a saved snapshot of a config file at a point in time.
type Entry struct {
	Label     string                 `json:"label"`
	File      string                 `json:"file"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Store manages snapshot persistence on disk.
type Store struct {
	Dir string
}

// NewStore creates a Store that saves snapshots under dir.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Store{Dir: dir}, nil
}

// Save writes a snapshot entry to disk using label as the filename key.
func (s *Store) Save(label, file string, data map[string]interface{}) error {
	entry := Entry{
		Label:     label,
		File:      file,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}
	b, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	dest := filepath.Join(s.Dir, label+".json")
	if err := os.WriteFile(dest, b, 0644); err != nil {
		return fmt.Errorf("snapshot: write: %w", err)
	}
	return nil
}

// Load reads a snapshot entry by label.
func (s *Store) Load(label string) (*Entry, error) {
	path := filepath.Join(s.Dir, label+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot %q not found", label)
		}
		return nil, fmt.Errorf("snapshot: read: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(b, &entry); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &entry, nil
}

// List returns all snapshot labels stored in the directory.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	var labels []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			labels = append(labels, e.Name()[:len(e.Name())-5])
		}
	}
	return labels, nil
}
