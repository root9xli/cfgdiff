package merge

import (
	"testing"
)

func baseMap() map[string]interface{} {
	return map[string]interface{}{
		"app.name":    "myapp",
		"app.port":    "8080",
		"app.debug":   "false",
		"db.host":     "localhost",
	}
}

func TestMerge_NoConflicts(t *testing.T) {
	base := baseMap()
	override := map[string]interface{}{
		"app.timeout": "30s",
	}

	res, err := Merge(base, override, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if res.Data["app.timeout"] != "30s" {
		t.Errorf("expected app.timeout=30s, got %v", res.Data["app.timeout"])
	}
	if res.Data["app.name"] != "myapp" {
		t.Errorf("base key app.name should be preserved")
	}
}

func TestMerge_StrategyOverride(t *testing.T) {
	base := baseMap()
	override := map[string]interface{}{
		"app.port": "9090",
	}

	res, err := Merge(base, override, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.Data["app.port"] != "9090" {
		t.Errorf("expected override value 9090, got %v", res.Data["app.port"])
	}
}

func TestMerge_StrategyBase(t *testing.T) {
	base := baseMap()
	override := map[string]interface{}{
		"app.port": "9090",
	}

	res, err := Merge(base, override, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Data["app.port"] != "8080" {
		t.Errorf("expected base value 8080, got %v", res.Data["app.port"])
	}
}

func TestMerge_StrategyError(t *testing.T) {
	base := baseMap()
	override := map[string]interface{}{
		"app.port": "9090",
	}

	_, err := Merge(base, override, StrategyError)
	if err == nil {
		t.Fatal("expected error for conflict with StrategyError, got nil")
	}
}

func TestMerge_IdenticalValues_NoConflict(t *testing.T) {
	base := baseMap()
	override := map[string]interface{}{
		"app.port": "8080", // same value
	}

	res, err := Merge(base, override, StrategyError)
	if err != nil {
		t.Fatalf("unexpected error for identical value: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts for identical values, got %d", len(res.Conflicts))
	}
}
