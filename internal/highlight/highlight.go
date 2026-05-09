// Package highlight provides terminal color highlighting for diff output.
package highlight

import (
	"fmt"
	"strings"

	"github.com/user/cfgdiff/internal/diff"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Options controls highlight behaviour.
type Options struct {
	Enabled bool // set false to disable all ANSI codes
}

// DefaultOptions returns Options with highlighting enabled.
func DefaultOptions() Options {
	return Options{Enabled: true}
}

// Highlighter applies ANSI colors to diff changes.
type Highlighter struct {
	opts Options
}

// New creates a Highlighter with the given options.
func New(opts Options) *Highlighter {
	return &Highlighter{opts: opts}
}

// Line returns a colorized single-line representation of a Change.
func (h *Highlighter) Line(c diff.Change) string {
	if !h.opts.Enabled {
		return h.plainLine(c)
	}
	switch c.Type {
	case diff.Added:
		return fmt.Sprintf("%s+ %s = %s%s", colorGreen, c.Key, c.NewValue, colorReset)
	case diff.Removed:
		return fmt.Sprintf("%s- %s = %s%s", colorRed, c.Key, c.OldValue, colorReset)
	case diff.Modified:
		return fmt.Sprintf("%s~ %s%s: %s%s%s → %s%s%s",
			colorYellow, colorBold, c.Key, colorReset,
			colorRed, c.OldValue, colorGreen, c.NewValue, colorReset)
	default:
		return h.plainLine(c)
	}
}

// Apply returns a slice of highlighted lines for all changes.
func (h *Highlighter) Apply(changes []diff.Change) []string {
	lines := make([]string, 0, len(changes))
	for _, c := range changes {
		lines = append(lines, h.Line(c))
	}
	return lines
}

// Header returns a bold, cyan-colored section header.
func (h *Highlighter) Header(text string) string {
	if !h.opts.Enabled {
		return text
	}
	return fmt.Sprintf("%s%s%s%s", colorCyan, colorBold, text, colorReset)
}

func (h *Highlighter) plainLine(c diff.Change) string {
	switch c.Type {
	case diff.Added:
		return fmt.Sprintf("+ %s = %s", c.Key, c.NewValue)
	case diff.Removed:
		return fmt.Sprintf("- %s = %s", c.Key, c.OldValue)
	case diff.Modified:
		return fmt.Sprintf("~ %s: %s → %s", c.Key, c.OldValue, c.NewValue)
	default:
		return strings.TrimSpace(c.Key)
	}
}
