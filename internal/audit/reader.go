package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Filter constrains which audit entries are returned.
type Filter struct {
	Since *time.Time
	Until *time.Time
	FileA string
	FileB string
}

// ReadLog reads all audit entries from the log file, applying an optional filter.
func ReadLog(path string, f *Filter) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open log: %w", err)
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		if f != nil && !matchesFilter(e, f) {
			continue
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("audit: scan log: %w", err)
	}
	return entries, nil
}

func matchesFilter(e Entry, f *Filter) bool {
	if f.Since != nil && e.Timestamp.Before(*f.Since) {
		return false
	}
	if f.Until != nil && e.Timestamp.After(*f.Until) {
		return false
	}
	if f.FileA != "" && e.FileA != f.FileA {
		return false
	}
	if f.FileB != "" && e.FileB != f.FileB {
		return false
	}
	return true
}
