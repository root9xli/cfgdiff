package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/cfgdiff/internal/snapshot"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	data := map[string]interface{}{"host": "localhost", "port": "5432"}
	if err := store.Save("prod", "config.yaml", data); err != nil {
		t.Fatalf("Save: %v", err)
	}

	entry, err := store.Load("prod")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if entry.Label != "prod" {
		t.Errorf("expected label prod, got %s", entry.Label)
	}
	if entry.Data["host"] != "localhost" {
		t.Errorf("expected host localhost, got %v", entry.Data["host"])
	}
}

func TestLoad_NotFound(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)
	_, err := store.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestList_Empty(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)
	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected empty list, got %v", labels)
	}
}

func TestList_MultipleSnapshots(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.NewStore(dir)

	data := map[string]interface{}{"key": "val"}
	_ = store.Save("alpha", "a.env", data)
	_ = store.Save("beta", "b.env", data)

	labels, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestNewStore_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "snaps")
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
