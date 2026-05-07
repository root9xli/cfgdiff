package schema

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintViolations writes a formatted table of violations to w.
func PrintViolations(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "schema validation passed — no violations found")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tMESSAGE")
	fmt.Fprintln(tw, "---\t-------")
	for _, v := range violations {
		fmt.Fprintf(tw, "%s\t%s\n", v.Key, v.Message)
	}
	_ = tw.Flush()
	fmt.Fprintf(w, "\n%d violation(s) found\n", len(violations))
}
