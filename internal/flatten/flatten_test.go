package flatten_test

import (
	"testing"

	"cfgdiff/internal/flatten"
)

func TestFlatten_SimpleMap(t *testing.T) {
	input := map[string]any{
		"host": "localhost",
		"port": 5432,
	}
	got := flatten.Flatten(input, flatten.DefaultOptions())
	if got["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %s", got["host"])
	}
	if got["port"] != "5432" {
		t.Errorf("expected port=5432, got %s", got["port"])
	}
}

func TestFlatten_NestedMap(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
	}
	got := flatten.Flatten(input, flatten.DefaultOptions())
	if got["db.host"] != "localhost" {
		t.Errorf("expected db.host=localhost, got %s", got["db.host"])
	}
	if got["db.port"] != "5432" {
		t.Errorf("expected db.port=5432, got %s", got["db.port"])
	}
	if _, ok := got["db"]; ok {
		t.Error("intermediate key 'db' should not appear in flat map")
	}
}

func TestFlatten_DeeplyNested(t *testing.T) {
	input := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	got := flatten.Flatten(input, flatten.DefaultOptions())
	if got["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %s", got["a.b.c"])
	}
}

func TestFlatten_CustomSeparator(t *testing.T) {
	input := map[string]any{
		"db": map[string]any{"host": "localhost"},
	}
	opts := flatten.Options{Separator: "__"}
	got := flatten.Flatten(input, opts)
	if got["db__host"] != "localhost" {
		t.Errorf("expected db__host=localhost, got %s", got["db__host"])
	}
}

func TestFlatten_WithPrefix(t *testing.T) {
	input := map[string]any{"key": "val"}
	opts := flatten.Options{Separator: ".", Prefix: "cfg"}
	got := flatten.Flatten(input, opts)
	if got["cfg.key"] != "val" {
		t.Errorf("expected cfg.key=val, got %s", got["cfg.key"])
	}
}

func TestFlatten_EmptyMap(t *testing.T) {
	got := flatten.Flatten(map[string]any{}, flatten.DefaultOptions())
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestKeys_Sorted(t *testing.T) {
	flat := map[string]string{"z": "1", "a": "2", "m": "3"}
	keys := flatten.Keys(flat)
	if keys[0] != "a" || keys[1] != "m" || keys[2] != "z" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestHasPrefix(t *testing.T) {
	if !flatten.HasPrefix("db.host", "db", ".") {
		t.Error("expected db.host to have prefix db")
	}
	if flatten.HasPrefix("dbname", "db", ".") {
		t.Error("dbname should not match prefix db with dot separator")
	}
	if !flatten.HasPrefix("db", "db", ".") {
		t.Error("exact match should return true")
	}
}
