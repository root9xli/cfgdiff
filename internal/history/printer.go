package history

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// PrintList writes a tabular summary of history entries to w.
func PrintList(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No history entries found.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTIMESTAMP\tFILE A\tFILE B\t+\t-\t~")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n",
			e.ID,
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.FileA,
			e.FileB,
			e.Summary.Added,
			e.Summary.Removed,
			e.Summary.Modified,
		)
	}
	tw.Flush()
}

// PrintEntry writes the full detail of a single history entry to w.
func PrintEntry(w io.Writer, e *Entry) {
	fmt.Fprintf(w, "ID:        %s\n", e.ID)
	fmt.Fprintf(w, "Timestamp: %s\n", e.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "File A:    %s\n", e.FileA)
	fmt.Fprintf(w, "File B:    %s\n", e.FileB)
	fmt.Fprintf(w, "Added:     %d\n", e.Summary.Added)
	fmt.Fprintf(w, "Removed:   %d\n", e.Summary.Removed)
	fmt.Fprintf(w, "Modified:  %d\n", e.Summary.Modified)
	if len(e.Changes) == 0 {
		fmt.Fprintln(w, "No changes recorded.")
		return
	}
	fmt.Fprintln(w, "\nChanges:")
	for _, c := range e.Changes {
		fmt.Fprintf(w, "  [%s] %s\n", c.Type, c.Key)
	}
}
