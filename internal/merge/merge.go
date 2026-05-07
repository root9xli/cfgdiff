package merge

import (
	"fmt"
	"maps"
)

// Strategy defines how conflicting keys are resolved.
type Strategy int

const (
	// StrategyBase keeps the base value on conflict.
	StrategyBase Strategy = iota
	// StrategyOverride replaces base value with override value on conflict.
	StrategyOverride
	// StrategyError returns an error on conflict.
	StrategyError
)

// Conflict records a key that existed in both maps with differing values.
type Conflict struct {
	Key      string
	BaseVal  interface{}
	OverVal  interface{}
}

// Result holds the merged flat map and any conflicts encountered.
type Result struct {
	Data      map[string]interface{}
	Conflicts []Conflict
}

// Merge combines base and override flat maps according to the given strategy.
// Both maps are expected to use dot-separated keys (as produced by parser.Parse).
func Merge(base, override map[string]interface{}, strategy Strategy) (*Result, error) {
	result := &Result{
		Data:      make(map[string]interface{}),
		Conflicts: []Conflict{},
	}

	// Copy base into result.
	maps.Copy(result.Data, base)

	for k, overVal := range override {
		baseVal, exists := result.Data[k]
		if !exists {
			result.Data[k] = overVal
			continue
		}

		if fmt.Sprintf("%v", baseVal) == fmt.Sprintf("%v", overVal) {
			// Values are identical — no conflict.
			continue
		}

		conflict := Conflict{Key: k, BaseVal: baseVal, OverVal: overVal}
		result.Conflicts = append(result.Conflicts, conflict)

		switch strategy {
		case StrategyBase:
			// Keep existing base value — do nothing.
		case StrategyOverride:
			result.Data[k] = overVal
		case StrategyError:
			return nil, fmt.Errorf("merge conflict on key %q: base=%v override=%v", k, baseVal, overVal)
		}
	}

	return result, nil
}
