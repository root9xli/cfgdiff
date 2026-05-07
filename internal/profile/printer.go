package profile

import (
	"fmt"
	"io"
	"strings"
)

// PrintList writes a human-readable table of profiles to w.
func PrintList(w io.Writer, profiles []Profile) {
	if len(profiles) == 0 {
		fmt.Fprintln(w, "no profiles saved")
		return
	}
	fmt.Fprintf(w, "%-20s  %s\n", "NAME", "FILES")
	fmt.Fprintf(w, "%-20s  %s\n", strings.Repeat("-", 20), strings.Repeat("-", 40))
	for _, p := range profiles {
		fmt.Fprintf(w, "%-20s  %s\n", p.Name, strings.Join(p.Files, ", "))
	}
}

// PrintProfile writes the details of a single profile to w.
func PrintProfile(w io.Writer, p Profile) {
	fmt.Fprintf(w, "Profile: %s\n", p.Name)
	if len(p.Files) == 0 {
		fmt.Fprintln(w, "  (no files)")
		return
	}
	for _, f := range p.Files {
		fmt.Fprintf(w, "  - %s\n", f)
	}
}
