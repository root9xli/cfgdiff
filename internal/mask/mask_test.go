package mask

import (
	"strings"
	"testing"
)

func TestNew_DefaultPatterns(t *testing.T) {
	m := New(DefaultOptions())
	sensitive := []string{"password", "db_password", "api_key", "API_KEY", "secret", "auth_token", "private_key"}
	for _, k := range sensitive {
		if !m.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_NoMatch(t *testing.T) {
	m := New(DefaultOptions())
	safe := []string{"host", "port", "database", "timeout", "region"}
	for _, k := range safe {
		if m.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestMaskValue_FullMask(t *testing.T) {
	m := New(DefaultOptions())
	result := m.MaskValue("supersecret")
	if result != "***********" {
		t.Errorf("expected all stars, got %q", result)
	}
}

func TestMaskValue_ShowFirst(t *testing.T) {
	m := New(Options{Char: "*", ShowFirst: 3, ShowLast: 0})
	result := m.MaskValue("abcdefgh")
	if !strings.HasPrefix(result, "abc") {
		t.Errorf("expected prefix 'abc', got %q", result)
	}
	if result != "abc*****" {
		t.Errorf("unexpected result %q", result)
	}
}

func TestMaskValue_ShowLast(t *testing.T) {
	m := New(Options{Char: "#", ShowFirst: 0, ShowLast: 2})
	result := m.MaskValue("abcdefgh")
	if result != "######gh" {
		t.Errorf("unexpected result %q", result)
	}
}

func TestMaskValue_ShowFirstAndLast(t *testing.T) {
	m := New(Options{Char: "*", ShowFirst: 2, ShowLast: 2})
	result := m.MaskValue("abcdefgh")
	if result != "ab****gh" {
		t.Errorf("unexpected result %q", result)
	}
}

func TestMaskValue_EmptyString(t *testing.T) {
	m := New(DefaultOptions())
	if got := m.MaskValue(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_MasksSensitiveKeys(t *testing.T) {
	m := New(DefaultOptions())
	input := map[string]string{
		"host":     "localhost",
		"password": "hunter2",
		"api_key":  "abc123",
	}
	out := m.Apply(input)
	if out["host"] != "localhost" {
		t.Errorf("host should be unchanged")
	}
	if out["password"] == "hunter2" {
		t.Errorf("password should be masked")
	}
	if out["api_key"] == "abc123" {
		t.Errorf("api_key should be masked")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m := New(DefaultOptions())
	input := map[string]string{"secret": "myvalue"}
	_ = m.Apply(input)
	if input["secret"] != "myvalue" {
		t.Errorf("original map was mutated")
	}
}

func TestNewWithPatterns_CustomPattern(t *testing.T) {
	m := NewWithPatterns(DefaultOptions(), []string{`(?i)internal`})
	if !m.IsSensitive("internal_url") {
		t.Errorf("expected internal_url to be sensitive")
	}
	if m.IsSensitive("password") {
		t.Errorf("default patterns should not apply with custom patterns")
	}
}
