package normalize

import (
	"testing"
)

func TestApply_TrimSpaceAndLowercase(t *testing.T) {
	input := map[string]string{
		"  APP_HOST  ": "  localhost  ",
		"APP_PORT":    "8080",
	}
	opts := DefaultOptions()
	result := Apply(input, opts)

	if v, ok := result["app_host"]; !ok || v != "localhost" {
		t.Errorf("expected app_host=localhost, got %q", v)
	}
	if v, ok := result["app_port"]; !ok || v != "8080" {
		t.Errorf("expected app_port=8080, got %q", v)
	}
}

func TestApply_StripPrefix(t *testing.T) {
	input := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "9090",
		"DB_NAME":  "mydb",
	}
	opts := Options{
		LowercaseKeys: false,
		TrimSpace:     false,
		StripPrefix:   "APP_",
	}
	result := Apply(input, opts)

	if _, ok := result["HOST"]; !ok {
		t.Error("expected key HOST after stripping APP_ prefix")
	}
	if _, ok := result["PORT"]; !ok {
		t.Error("expected key PORT after stripping APP_ prefix")
	}
	if _, ok := result["DB_NAME"]; !ok {
		t.Error("expected key DB_NAME to remain unchanged")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := map[string]string{"KEY": "value"}
	opts := DefaultOptions()
	_ = Apply(input, opts)

	if _, ok := input["KEY"]; !ok {
		t.Error("original map should not be mutated")
	}
	if _, ok := input["key"]; ok {
		t.Error("original map should not have lowercased key")
	}
}

func TestApply_NoOptions(t *testing.T) {
	input := map[string]string{"  KEY  ": "  val  "}
	opts := Options{}
	result := Apply(input, opts)

	if v, ok := result["  KEY  "]; !ok || v != "  val  " {
		t.Errorf("expected keys/values to be unchanged, got %q", v)
	}
}

func TestNormalizeKey_Standalone(t *testing.T) {
	opts := Options{LowercaseKeys: true, TrimSpace: true, StripPrefix: "ENV_"}
	got := NormalizeKey("  ENV_DATABASE  ", opts)
	want := "database"
	if got != want {
		t.Errorf("NormalizeKey: got %q, want %q", got, want)
	}
}
