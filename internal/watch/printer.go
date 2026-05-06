package watch

import (
	"fmt"
	"io"
	"time"
)

// PrintChange writes a human-readable change notification to w.
func PrintChange(w io.Writer, path string, prev, curr FileState) {
	timestamp := curr.ModTime.Format(time.RFC3339)
	if prev.Checksum == "" {
		fmt.Fprintf(w, "[%s] DETECTED  %s (first seen, checksum: %s)\n",
			timestamp, path, curr.Checksum)
		return
	}
	fmt.Fprintf(w, "[%s] CHANGED   %s\n", timestamp, path)
	fmt.Fprintf(w, "  prev checksum: %s\n", prev.Checksum)
	fmt.Fprintf(w, "  curr checksum: %s\n", curr.Checksum)
}
