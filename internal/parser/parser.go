package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFormat represents the supported config file formats.
type ConfigFormat string

const (
	FormatJSON ConfigFormat = "json"
	FormatYAML ConfigFormat = "yaml"
	FormatTOML ConfigFormat = "toml"
	FormatENV  ConfigFormat = "env"
)

// ConfigFile holds the parsed representation of a config file.
type ConfigFile struct {
	Path   string
	Format ConfigFormat
	Data   map[string]interface{}
}

// DetectFormat infers the config format from the file extension.
func DetectFormat(path string) (ConfigFormat, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	switch ext {
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	case "toml":
		return FormatTOML, nil
	case "env":
		return FormatENV, nil
	default:
		return "", fmt.Errorf("unsupported file extension: %q", ext)
	}
}

// Parse reads and parses a config file, returning a ConfigFile.
func Parse(path string) (*ConfigFile, error) {
	format, err := DetectFormat(path)
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	var data map[string]interface{}

	switch format {
	case FormatJSON:
		data, err = parseJSON(raw)
	case FormatYAML:
		data, err = parseYAML(raw)
	case FormatTOML:
		data, err = parseTOML(raw)
	case FormatENV:
		data, err = parseENV(raw)
	}

	if err != nil {
		return nil, fmt.Errorf("parsing %s file %q: %w", format, path, err)
	}

	return &ConfigFile{
		Path:   path,
		Format: format,
		Data:   data,
	}, nil
}
