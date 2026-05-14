package patch

import (
	"fmt"
	"io"
	"os"

	"github.com/cfgdiff/internal/diff"
)

// PrintPatch writes a human-readable summary of the patch to stdout.
func PrintPatch(p *Patch) {
	FprintPatch(os.Stdout, p)
}

// FprintPatch writes a human-readable summary of the patch to w.
func FprintPatch(w io.Writer, p *Patch) {
	dir := "forward"
	if p.Direction == Reverse {
		dir = "reverse"
	}
	fmt.Fprintf(w, "Patch (%s) — %d change(s)\n", dir, len(p.Changes))
	fmt.Fprintln(w, "---")
	for _, c := range p.Changes {
		switch c.Type {
		case diff.Added:
			fmt.Fprintf(w, "  + %-30s %s\n", c.Key, c.NewValue)
		case diff.Removed:
			fmt.Fprintf(w, "  - %-30s %s\n", c.Key, c.OldValue)
		case diff.Modified:
			fmt.Fprintf(w, "  ~ %-30s %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		}
	}
}
