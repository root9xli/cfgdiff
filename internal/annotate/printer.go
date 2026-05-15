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
// If there are no annotations, it prints a message indicating so.
func Fprint(w io.Writer, annotations []Annotation) error {
	if len(annotations) == 0 {
		_, err := fmt.Fprintln(w, "no annotations")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tKEY\tNOTE")
	fmt.Fprintln(tw, "----\t---\t----")
	for _, a := range annotations {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", a.Change.Type, a.Change.Key, a.Note)
	}
	return tw.Flush()
}
