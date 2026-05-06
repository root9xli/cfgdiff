package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestRootCmd_NoDiff(t *testing.T) {
	a := writeTempConfig(t, "a.json", `{"key":"value"}`)
	b := writeTempConfig(t, "b.json", `{"key":"value"}`)

	rootCmd.SetArgs([]string{a, b})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCmd_WithChanges(t *testing.T) {
	a := writeTempConfig(t, "a.yaml", "host: localhost\nport: 8080\n")
	b := writeTempConfig(t, "b.yaml", "host: prod.example.com\nport: 443\n")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{a, b})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRootCmd_InvalidFile(t *testing.T) {
	rootCmd.SetArgs([]string{"/nonexistent/a.json", "/nonexistent/b.json"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing files, got nil")
	}
}

func TestRootCmd_JSONFormat(t *testing.T) {
	a := writeTempConfig(t, "a.json", `{"env":"dev"}`)
	b := writeTempConfig(t, "b.json", `{"env":"prod"}`)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--format", "json", a, b})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "{") {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestAuditCmd_MissingFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"audit"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --audit-log not provided")
	}
}

func TestAuditCmd_WithLog(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "audit.log")
	// write a minimal valid JSONL entry
	entry := `{"timestamp":"2024-01-01T00:00:00Z","file_a":"a.json","file_b":"b.json","added":0,"removed":0,"modified":0}` + "\n"
	if err := os.WriteFile(logFile, []byte(entry), 0644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"--audit-log", logFile, "audit"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
