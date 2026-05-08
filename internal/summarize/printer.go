package summarize

import (
	"fmt"
	"io"
	"strings"
)

// Print writes a formatted summary block to w.
func Print(w io.Writer, s Stats) {
	fmt.Fprintln(w, "=== Change Summary ===")
	fmt.Fprintf(w, "  Total    : %d\n", s.Total)
	fmt.Fprintf(w, "  Added    : %d\n", s.Added)
	fmt.Fprintf(w, "  Removed  : %d\n", s.Removed)
	fmt.Fprintf(w, "  Modified : %d\n", s.Modified)

	if len(s.TopKeys) > 0 {
		fmt.Fprintf(w, "  Top keys : %s\n", strings.Join(s.TopKeys, ", "))
	}
}
