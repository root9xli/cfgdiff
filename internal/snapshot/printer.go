package snapshot

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PrintList writes a formatted table of snapshot labels to w.
func PrintList(w io.Writer, labels []string) {
	if len(labels) == 0 {
		fmt.Fprintln(w, "No snapshots saved.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "LABEL")
	sort.Strings(labels)
	for _, l := range labels {
		fmt.Fprintln(tw, l)
	}
	tw.Flush()
}

// PrintEntry writes the details of a single snapshot entry to w.
func PrintEntry(w io.Writer, e *Entry) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Label:\t%s\n", e.Label)
	fmt.Fprintf(tw, "File:\t%s\n", e.File)
	fmt.Fprintf(tw, "Timestamp:\t%s\n", e.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(tw, "Keys:\t%d\n", len(e.Data))
	tw.Flush()
}
