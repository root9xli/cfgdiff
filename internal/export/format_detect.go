package export

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DetectFormat infers the export Format from a file extension.
// It returns an error if the extension is unrecognised.
func DetectFormat(filename string) (Format, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	switch ext {
	case "csv":
		return FormatCSV, nil
	case "md", "markdown":
		return FormatMarkdown, nil
	case "json":
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("cannot detect export format from extension %q", ext)
	}
}
