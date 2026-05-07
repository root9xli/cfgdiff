package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/cfgdiff/internal/diff"
)

// Format represents a supported export format.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
	FormatJSON     Format = "json"
)

// Exporter writes diff changes to an output stream in a given format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New creates a new Exporter for the given format and writer.
func New(format Format, w io.Writer) (*Exporter, error) {
	switch format {
	case FormatCSV, FormatMarkdown, FormatJSON:
		return &Exporter{format: format, w: w}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// Write exports the given changes to the configured writer.
func (e *Exporter) Write(changes []diff.Change) error {
	switch e.format {
	case FormatCSV:
		return writeCSV(e.w, changes)
	case FormatMarkdown:
		return writeMarkdown(e.w, changes)
	case FormatJSON:
		return writeJSON(e.w, changes)
	default:
		return fmt.Errorf("unsupported format: %s", e.format)
	}
}

func writeCSV(w io.Writer, changes []diff.Change) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "type", "old_value", "new_value"}); err != nil {
		return err
	}
	for _, c := range changes {
		row := []string{c.Key, string(c.Type), fmt.Sprintf("%v", c.OldValue), fmt.Sprintf("%v", c.NewValue)}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(w io.Writer, changes []diff.Change) error {
	lines := []string{
		"| Key | Type | Old Value | New Value |",
		"|-----|------|-----------|-----------|\n",
	}
	for _, c := range changes {
		row := fmt.Sprintf("| %s | %s | %v | %v |", c.Key, c.Type, c.OldValue, c.NewValue)
		lines = append(lines, row)
	}
	_, err := fmt.Fprintln(w, strings.Join(lines, "\n"))
	return err
}

func writeJSON(w io.Writer, changes []diff.Change) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(changes)
}
