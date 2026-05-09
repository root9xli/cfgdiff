package classify

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintResults writes classified results to stdout in a table.
func PrintResults(results []Result) {
	FprintResults(os.Stdout, results)
}

// FprintResults writes classified results to w.
func FprintResults(w io.Writer, results []Result) {
	if len(results) == 0 {
		fmt.Fprintln(w, "no changes to classify")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "KEY\tTYPE\tDOMAIN\tSEVERITY")
	fmt.Fprintln(tw, "---\t----\t------\t--------")
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			r.Change.Key,
			r.Change.Type,
			r.Domain,
			r.Severity,
		)
	}
	tw.Flush()
}
