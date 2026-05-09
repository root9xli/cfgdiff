package highlight

import (
	"fmt"
	"io"
	"os"

	"github.com/user/cfgdiff/internal/diff"
)

// Print writes highlighted diff lines to stdout.
func Print(changes []diff.Change, opts Options) {
	Fprint(os.Stdout, changes, opts)
}

// Fprint writes highlighted diff lines to the given writer.
func Fprint(w io.Writer, changes []diff.Change, opts Options) {
	h := New(opts)
	if len(changes) == 0 {
		fmt.Fprintln(w, h.Header("No differences found."))
		return
	}
	fmt.Fprintln(w, h.Header(fmt.Sprintf("=== %d change(s) ===", len(changes))))
	for _, line := range h.Apply(changes) {
		fmt.Fprintln(w, line)
	}
}
