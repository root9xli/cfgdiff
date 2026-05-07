package schema

import (
	"testing"
)

func TestValidate_AllValid(t *testing.T) {
	s := Load(map[string]FieldRule{
		"HOST": {Required: true, Type: "string"},
		"PORT": {Required: true, Type: "int"},
		"DEBUG": {Type: "bool"},
	})
	cfg := map[string]string{"HOST": "localhost", "PORT": "8080", "DEBUG": "true"}
	if v := s.Validate(cfg); len(v) != 0 {
		t.Errorf("expected no violations, got %v", v)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := Load(map[string]FieldRule{
		"DB_URL": {Required: true},
	})
	v := s.Validate(map[string]string{})
	if len(v) != 1 || v[0].Key != "DB_URL" {
		t.Errorf("expected DB_URL violation, got %v", v)
	}
}

func TestValidate_WrongType_Int(t *testing.T) {
	s := Load(map[string]FieldRule{
		"PORT": {Type: "int"},
	})
	v := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_WrongType_Bool(t *testing.T) {
	s := Load(map[string]FieldRule{
		"ENABLED": {Type: "bool"},
	})
	v := s.Validate(map[string]string{"ENABLED": "yes"})
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	s := Load(map[string]FieldRule{
		"ENV": {Pattern: `^(production|staging|development)$`},
	})
	v := s.Validate(map[string]string{"ENV": "prod"})
	if len(v) != 1 {
		t.Errorf("expected pattern violation, got %v", v)
	}
}

func TestValidate_AllowedValues(t *testing.T) {
	s := Load(map[string]FieldRule{
		"LOG_LEVEL": {Allowed: []string{"debug", "info", "warn", "error"}},
	})
	v := s.Validate(map[string]string{"LOG_LEVEL": "verbose"})
	if len(v) != 1 {
		t.Errorf("expected allowed-values violation, got %v", v)
	}
	v2 := s.Validate(map[string]string{"LOG_LEVEL": "info"})
	if len(v2) != 0 {
		t.Errorf("expected no violations for valid allowed value")
	}
}

func TestViolation_String(t *testing.T) {
	v := Violation{Key: "PORT", Message: "expected int"}
	if v.String() != "PORT: expected int" {
		t.Errorf("unexpected string: %s", v.String())
	}
}
