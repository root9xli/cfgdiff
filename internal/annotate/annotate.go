// Package annotate attaches human-readable comments to diff changes.
package annotate

import (
	"fmt"
	"strings"

	"github.com/user/cfgdiff/internal/diff"
)

// Annotation holds a change along with a descriptive note.
type Annotation struct {
	Change diff.Change
	Note   string
}

// Options controls annotation behaviour.
type Options struct {
	// CustomNotes maps key prefixes to a static note string.
	CustomNotes map[string]string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		CustomNotes: map[string]string{},
	}
}

// Annotator applies notes to a slice of changes.
type Annotator struct {
	opts Options
}

// New creates an Annotator with the given options.
func New(opts Options) *Annotator {
	return &Annotator{opts: opts}
}

// Apply returns a slice of Annotations for the provided changes.
func (a *Annotator) Apply(changes []diff.Change) []Annotation {
	out := make([]Annotation, 0, len(changes))
	for _, c := range changes {
		out = append(out, Annotation{
			Change: c,
			Note:   a.noteFor(c),
		})
	}
	return out
}

// noteFor derives a note string for a single change.
func (a *Annotator) noteFor(c diff.Change) string {
	// Check user-supplied prefix notes first.
	for prefix, note := range a.opts.CustomNotes {
		if strings.HasPrefix(c.Key, prefix) {
			return note
		}
	}
	switch c.Type {
	case diff.Added:
		return fmt.Sprintf("key %q introduced with value %q", c.Key, c.NewValue)
	case diff.Removed:
		return fmt.Sprintf("key %q removed (was %q)", c.Key, c.OldValue)
	case diff.Modified:
		return fmt.Sprintf("key %q changed from %q to %q", c.Key, c.OldValue, c.NewValue)
	default:
		return ""
	}
}
