// Package profile manages named environment profiles, allowing users to
// store and retrieve sets of config file paths grouped under a label.
package profile

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Profile represents a named collection of config file paths.
type Profile struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// Store manages profiles persisted to a JSON file on disk.
type Store struct {
	path string
}

// NewStore creates a Store backed by the given directory.
// The directory is created if it does not exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(dir, "profiles.json")}, nil
}

// Save adds or replaces a profile in the store.
func (s *Store) Save(p Profile) error {
	profiles, err := s.loadAll()
	if err != nil {
		return err
	}
	profiles[p.Name] = p
	return s.writeAll(profiles)
}

// Get retrieves a profile by name.
func (s *Store) Get(name string) (Profile, error) {
	profiles, err := s.loadAll()
	if err != nil {
		return Profile{}, err
	}
	p, ok := profiles[name]
	if !ok {
		return Profile{}, errors.New("profile not found: " + name)
	}
	return p, nil
}

// Delete removes a profile by name. Returns an error if it does not exist.
func (s *Store) Delete(name string) error {
	profiles, err := s.loadAll()
	if err != nil {
		return err
	}
	if _, ok := profiles[name]; !ok {
		return errors.New("profile not found: " + name)
	}
	delete(profiles, name)
	return s.writeAll(profiles)
}

// List returns all stored profiles.
func (s *Store) List() ([]Profile, error) {
	profiles, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	out := make([]Profile, 0, len(profiles))
	for _, p := range profiles {
		out = append(out, p)
	}
	return out, nil
}

func (s *Store) loadAll() (map[string]Profile, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]Profile{}, nil
	}
	if err != nil {
		return nil, err
	}
	var profiles map[string]Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (s *Store) writeAll(profiles map[string]Profile) error {
	data, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
