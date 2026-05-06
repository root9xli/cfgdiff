package ignore

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "cfgignore-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestNewRules_NoMatch(t *testing.T) {
	r := NewRules([]string{})
	if r.Match("any.key") {
		t.Error("empty rules should not match any key")
	}
}

func TestMatch_ExactPattern(t *testing.T) {
	r := NewRules([]string{"DB_PASSWORD", "API_SECRET"})
	if !r.Match("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to match")
	}
	if r.Match("DB_HOST") {
		t.Error("DB_HOST should not match")
	}
}

func TestMatch_WildcardPattern(t *testing.T) {
	r := NewRules([]string{"secret.*"})
	if !r.Match("secret.key") {
		t.Error("expected secret.key to match secret.*")
	}
	if !r.Match("secret.token") {
		t.Error("expected secret.token to match secret.*")
	}
	if r.Match("public.key") {
		t.Error("public.key should not match secret.*")
	}
}

func TestMatch_GlobPattern(t *testing.T) {
	r := NewRules([]string{"*_PASSWORD"})
	if !r.Match("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to match *_PASSWORD")
	}
	if !r.Match("APP_PASSWORD") {
		t.Error("expected APP_PASSWORD to match *_PASSWORD")
	}
}

func TestFilterKeys(t *testing.T) {
	r := NewRules([]string{"DB_PASSWORD", "*_SECRET"})
	keys := []string{"DB_HOST", "DB_PASSWORD", "APP_SECRET", "APP_PORT"}
	filtered := r.FilterKeys(keys)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(filtered), filtered)
	}
	if filtered[0] != "DB_HOST" || filtered[1] != "APP_PORT" {
		t.Errorf("unexpected filtered keys: %v", filtered)
	}
}

func TestLoadFile_NotExist(t *testing.T) {
	r, err := LoadFile("/nonexistent/path/.cfgignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if r.Match("anything") {
		t.Error("rules from missing file should match nothing")
	}
}

func TestLoadFile_WithPatterns(t *testing.T) {
	path := writeTempIgnore(t, "# comment\n\nDB_PASSWORD\n*_TOKEN\n")
	r, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Match("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to match")
	}
	if !r.Match("GITHUB_TOKEN") {
		t.Error("expected GITHUB_TOKEN to match *_TOKEN")
	}
	if r.Match("DB_HOST") {
		t.Error("DB_HOST should not match")
	}
}
