package audit

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintEntries writes a human-readable audit log table to w.
func PrintEntries(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No audit entries found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tFILE_A\tFILE_B\tADDED\tREMOVED\tMODIFIED\tTOTAL")
	fmt.Fprintln(tw, "---------\t------\t------\t-----\t-------\t--------\t-----")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%d\t%d\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.FileA,
			e.FileB,
			e.Summary.Added,
			e.Summary.Removed,
			e.Summary.Modified,
			e.Summary.Total,
		)
	}
	tw.Flush()
}
