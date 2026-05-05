package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return path
}

func TestDetectFormat(t *testing.T) {
	cases := []struct {
		path   string
		want   ConfigFormat
		wantErr bool
	}{
		{"config.json", FormatJSON, false},
		{"config.yaml", FormatYAML, false},
		{"config.yml", FormatYAML, false},
		{"config.toml", FormatTOML, false},
		{"config.env", FormatENV, false},
		{"config.xml", "", true},
	}
	for _, tc := range cases {
		got, err := DetectFormat(tc.path)
		if tc.wantErr && err == nil {
			t.Errorf("DetectFormat(%q): expected error, got nil", tc.path)
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("DetectFormat(%q) = %q, want %q", tc.path, got, tc.want)
		}
	}
}

func TestParseJSON(t *testing.T) {
	path := writeTempFile(t, "config.json", `{"host": "localhost", "port": 8080}`)
	cf, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cf.Data["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cf.Data["host"])
	}
}

func TestParseYAML(t *testing.T) {
	path := writeTempFile(t, "config.yaml", "host: localhost\nport: 8080\n")
	cf, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cf.Data["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", cf.Data["host"])
	}
}

func TestParseENV(t *testing.T) {
	path := writeTempFile(t, "config.env", "# comment\nDB_HOST=localhost\nDB_PASS=\"secret\"\n")
	cf, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if cf.Data["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", cf.Data["DB_HOST"])
	}
	if cf.Data["DB_PASS"] != "secret" {
		t.Errorf("expected DB_PASS=secret, got %v", cf.Data["DB_PASS"])
	}
}

func TestParseUnsupportedFormat(t *testing.T) {
	path := writeTempFile(t, "config.xml", "<root/>")
	_, err := Parse(path)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
