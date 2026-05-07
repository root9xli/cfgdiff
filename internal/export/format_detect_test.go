package export

import (
	"testing"
)

func TestDetectFormat_CSV(t *testing.T) {
	f, err := DetectFormat("report.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != FormatCSV {
		t.Errorf("expected csv, got %s", f)
	}
}

func TestDetectFormat_Markdown(t *testing.T) {
	for _, name := range []string{"report.md", "report.markdown"} {
		f, err := DetectFormat(name)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", name, err)
		}
		if f != FormatMarkdown {
			t.Errorf("expected markdown for %s, got %s", name, f)
		}
	}
}

func TestDetectFormat_JSON(t *testing.T) {
	f, err := DetectFormat("output.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != FormatJSON {
		t.Errorf("expected json, got %s", f)
	}
}

func TestDetectFormat_Unknown(t *testing.T) {
	_, err := DetectFormat("report.xml")
	if err == nil {
		t.Error("expected error for unknown extension")
	}
}

func TestDetectFormat_CaseInsensitive(t *testing.T) {
	f, err := DetectFormat("REPORT.CSV")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != FormatCSV {
		t.Errorf("expected csv, got %s", f)
	}
}
