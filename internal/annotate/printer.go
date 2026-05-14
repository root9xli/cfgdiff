package annotate

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

// Print writes annotations to stdout.
func Print(annotations []Annotation) {
	Fprint(os.Stdout, annotations)
}

// Fprint writes annotations to the given writer in a tabular format.
func Fprint(w io.Writer, annotations []Annotation) {
	if len(annotations) == 0 {
		fmt.Fprintln(w, "no annotations")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tKEY\tNOTE")
	fmt.Fprintln(tw, "----\t---\t----")
	for _, a := range annotations {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", a.Change.Type, a.Change.Key, a.Note)
	}
	tw.Flush()
}
