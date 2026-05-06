package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cfgdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time    `json:"timestamp"`
	FileA     string       `json:"file_a"`
	FileB     string       `json:"file_b"`
	Changes   []diff.Change `json:"changes"`
	Summary   Summary      `json:"summary"`
}

// Summary holds counts of change types.
type Summary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

// Logger writes audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Record writes an audit entry for a diff operation.
func (l *Logger) Record(fileA, fileB string, changes []diff.Change) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		FileA:     fileA,
		FileB:     fileB,
		Changes:   changes,
		Summary:   buildSummary(changes),
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// buildSummary computes change type counts.
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
	s.Total = s.Added + s.Removed + s.Modified
	return s
}
