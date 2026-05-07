package validate

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintViolations writes a formatted table of violations to w.
// If there are no violations, a clean message is printed instead.
func PrintViolations(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "✔  no validation violations found")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tRULE\tMESSAGE")
	fmt.Fprintln(tw, "---\t----\t-------")
	for _, v := range violations {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", v.Key, v.Rule, v.Message)
	}
	tw.Flush()
	fmt.Fprintf(w, "\n%d violation(s) found\n", len(violations))
}
