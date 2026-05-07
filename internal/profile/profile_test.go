package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s
}

func TestSave_AndGet(t *testing.T) {
	s := newTestStore(t)
	p := Profile{Name: "staging", Files: []string{"a.yaml", "b.env"}}
	if err := s.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := s.Get("staging")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != p.Name || len(got.Files) != 2 {
		t.Errorf("unexpected profile: %+v", got)
	}
}

func TestGet_NotFound(t *testing.T) {
	s := newTestStore(t)
	_, err := s.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestDelete_RemovesProfile(t *testing.T) {
	s := newTestStore(t)
	_ = s.Save(Profile{Name: "prod", Files: []string{"prod.toml"}})
	if err := s.Delete("prod"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("prod")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := newTestStore(t)
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error deleting non-existent profile")
	}
}

func TestList_Empty(t *testing.T) {
	s := newTestStore(t)
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %d", len(list))
	}
}

func TestList_MultipleProfiles(t *testing.T) {
	s := newTestStore(t)
	_ = s.Save(Profile{Name: "dev", Files: []string{"dev.yaml"}})
	_ = s.Save(Profile{Name: "prod", Files: []string{"prod.yaml"}})
	list, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(list))
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "profiles")
	_, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to exist: %v", err)
	}
}
