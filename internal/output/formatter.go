package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"

	"cfgdiff/internal/diff"
)

// Format represents an output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes diff results to an output writer.
type Formatter struct {
	Writer io.Writer
	Format Format
	NoColor bool
}

// NewFormatter creates a Formatter with the given writer and format.
func NewFormatter(w io.Writer, format Format, noColor bool) *Formatter {
	if noColor {
		color.NoColor = true
	}
	return &Formatter{Writer: w, Format: format, NoColor: noColor}
}

// Write outputs the diff results according to the configured format.
func (f *Formatter) Write(results []diff.Change) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

func (f *Formatter) writeText(results []diff.Change) error {
	if len(results) == 0 {
		fmt.Fprintln(f.Writer, "No differences found.")
		return nil
	}

	added := color.New(color.FgGreen)
	removed := color.New(color.FgRed)
	modified := color.New(color.FgYellow)

	for _, c := range results {
		switch c.Type {
		case diff.Added:
			added.Fprintf(f.Writer, "+ [%s] = %v\n", c.Key, c.NewValue)
		case diff.Removed:
			removed.Fprintf(f.Writer, "- [%s] = %v\n", c.Key, c.OldValue)
		case diff.Modified:
			modified.Fprintf(f.Writer, "~ [%s]: %v → %v\n", c.Key, c.OldValue, c.NewValue)
		}
	}

	fmt.Fprintf(f.Writer, "\nSummary: %s\n", buildSummary(results))
	return nil
}

func (f *Formatter) writeJSON(results []diff.Change) error {
	// Simple JSON serialization without extra dependencies.
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, c := range results {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"type\": %q, \"old\": %q, \"new\": %q}",
			c.Key, c.Type, fmt.Sprintf("%v", c.OldValue), fmt.Sprintf("%v", c.NewValue),
		))
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(f.Writer, sb.String())
	return err
}

func buildSummary(results []diff.Change) string {
	var added, removed, modified int
	for _, c := range results {
		switch c.Type {
		case diff.Added:
			added++
		case diff.Removed:
			removed++
		case diff.Modified:
			modified++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d modified", added, removed, modified)
}
