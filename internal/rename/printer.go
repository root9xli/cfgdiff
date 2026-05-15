package rename

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// PrintRules prints the active rename rules in a human-readable table.
func PrintRules(rules []Rule) {
	FprintRules(os.Stdout, rules)
}

// FprintRules writes the rename rule table to the given writer.
func FprintRules(w io.Writer, rules []Rule) {
	if len(rules) == 0 {
		fmt.Fprintln(w, "No rename rules defined.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "FROM\tTO")
	fmt.Fprintln(tw, "----\t--")
	for _, r := range rules {
		fmt.Fprintf(tw, "%s\t%s\n", r.From, r.To)
	}
	tw.Flush()
}
