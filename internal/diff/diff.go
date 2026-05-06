package diff

import (
	"fmt"
	"sort"
)

// ChangeType represents the type of change between two configs.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change represents a single key-level difference between two config maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Result holds the full diff output between two configs.
type Result struct {
	Changes []Change
}

// HasChanges returns true if the diff contains any changes.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a human-readable summary of the diff.
func (r *Result) Summary() string {
	if !r.HasChanges() {
		return "No differences found."
	}
	added, removed, modified := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d modified", added, removed, modified)
}

// Compare performs a flat key-by-key diff between two parsed config maps.
func Compare(base, target map[string]interface{}) *Result {
	result := &Result{}

	for key, baseVal := range base {
		if targetVal, ok := target[key]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: baseVal,
				NewValue: nil,
			})
		} else if fmt.Sprintf("%v", baseVal) != fmt.Sprintf("%v", targetVal) {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Modified,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Type:     Added,
				OldValue: nil,
				NewValue: targetVal,
			})
		}
	}

	sort.Slice(result.Changes, func(i, j int) bool {
		return result.Changes[i].Key < result.Changes[j].Key
	})

	return result
}
