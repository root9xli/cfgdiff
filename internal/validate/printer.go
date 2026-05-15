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

// PrintViolationsByRule writes a formatted table of violations grouped by rule to w.
// Each rule is printed as a section header followed by its associated violations.
func PrintViolationsByRule(w io.Writer, violations []Violation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "✔  no validation violations found")
		return
	}

	// Group violations by rule name.
	groups := make(map[string][]Violation)
	order := []string{}
	for _, v := range violations {
		if _, seen := groups[v.Rule]; !seen {
			order = append(order, v.Rule)
		}
		groups[v.Rule] = append(groups[v.Rule], v)
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, rule := range order {
		fmt.Fprintf(tw, "[%s]\n", rule)
		for _, v := range groups[rule] {
			fmt.Fprintf(tw, "  %s\t%s\n", v.Key, v.Message)
		}
		fmt.Fprintln(tw)
	}
	tw.Flush()
	fmt.Fprintf(w, "%d violation(s) found\n", len(violations))
}
