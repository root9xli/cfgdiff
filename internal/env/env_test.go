package env

import (
	"testing"
)

func staticResolver() *Resolver {
	return NewWithMap(map[string]string{
		"HOST":    "localhost",
		"PORT":    "5432",
		"DB_NAME": "mydb",
	})
}

func TestResolve_NoReferences(t *testing.T) {
	r := staticResolver()
	input := map[string]any{"key": "plain-value", "num": 42}
	out := r.Resolve(input)
	if out["key"] != "plain-value" {
		t.Errorf("expected plain-value, got %v", out["key"])
	}
	if out["num"] != 42 {
		t.Errorf("expected 42, got %v", out["num"])
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	r := staticResolver()
	out := r.Resolve(map[string]any{"dsn": "${HOST}:${PORT}/${DB_NAME}"})
	want := "localhost:5432/mydb"
	if out["dsn"] != want {
		t.Errorf("expected %q, got %q", want, out["dsn"])
	}
}

func TestResolve_DollarStyle(t *testing.T) {
	r := staticResolver()
	out := r.Resolve(map[string]any{"host": "$HOST"})
	if out["host"] != "localhost" {
		t.Errorf("expected localhost, got %v", out["host"])
	}
}

func TestResolve_UnknownVarLeftAsIs(t *testing.T) {
	r := staticResolver()
	out := r.Resolve(map[string]any{"val": "${UNDEFINED_VAR}"})
	if out["val"] != "${UNDEFINED_VAR}" {
		t.Errorf("expected reference to remain, got %v", out["val"])
	}
}

func TestResolve_NonStringUnchanged(t *testing.T) {
	r := staticResolver()
	out := r.Resolve(map[string]any{"flag": true, "count": 7})
	if out["flag"] != true {
		t.Errorf("expected true, got %v", out["flag"])
	}
	if out["count"] != 7 {
		t.Errorf("expected 7, got %v", out["count"])
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	r := staticResolver()
	input := map[string]any{"addr": "${HOST}:${PORT}"}
	_ = r.Resolve(input)
	if input["addr"] != "${HOST}:${PORT}" {
		t.Error("Resolve must not mutate the original map")
	}
}

func TestNew_UsesOSEnv(t *testing.T) {
	// Smoke-test that New() constructs without panic and returns a non-nil resolver.
	r := New()
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}
