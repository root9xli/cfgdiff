package lint

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintViolations writes lint violations to w in a human-readable table.
// If there are no violations, a clean message is printed instead.
func PrintViolations(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "✔  no lint violations found")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tRULE\tMESSAGE")
	fmt.Fprintln(tw, "---\t----\t-------")
	for _, v := range violations {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", v.Key, v.Rule, v.Message)
	}
	_ = tw.Flush()
	fmt.Fprintf(w, "\n%d violation(s) found\n", len(violations))
}
