package redact

import (
	"testing"
)

func TestNew_DefaultPatterns(t *testing.T) {
	r, err := New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
	if r.mask != "***REDACTED***" {
		t.Errorf("expected default mask, got %q", r.mask)
	}
}

func TestIsSensitive_Matches(t *testing.T) {
	r, _ := New(nil, "")
	cases := []struct {
		key       string
		expected  bool
	}{
		{"password", true},
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"MY_SECRET_TOKEN", true},
		{"database_host", false},
		{"port", false},
		{"username", false},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := r.IsSensitive(tc.key)
			if got != tc.expected {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.expected)
			}
		})
	}
}

func TestApply_RedactsValues(t *testing.T) {
	r, _ := New(nil, "[hidden]")
	input := map[string]any{
		"host":     "localhost",
		"port":     5432,
		"password": "supersecret",
		"api_key":  "abc123",
	}
	result := r.Apply(input)

	if result["host"] != "localhost" {
		t.Errorf("expected host unchanged, got %v", result["host"])
	}
	if result["port"] != 5432 {
		t.Errorf("expected port unchanged, got %v", result["port"])
	}
	if result["password"] != "[hidden]" {
		t.Errorf("expected password redacted, got %v", result["password"])
	}
	if result["api_key"] != "[hidden]" {
		t.Errorf("expected api_key redacted, got %v", result["api_key"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r, _ := New(nil, "")
	input := map[string]any{
		"password": "original",
	}
	_ = r.Apply(input)
	if input["password"] != "original" {
		t.Error("Apply mutated the original map")
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r, err := New([]string{"pin", "cvv"}, "XXXX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsSensitive("pin_code") {
		t.Error("expected pin_code to be sensitive")
	}
	if !r.IsSensitive("CVV") {
		t.Error("expected CVV to be sensitive")
	}
	if r.IsSensitive("password") {
		t.Error("expected password NOT sensitive with custom patterns")
	}
}
