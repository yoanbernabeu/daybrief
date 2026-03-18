package newsletter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "newsletter.json")

	content := `{
		"generated_at": "2026-03-18T09:00:00Z",
		"subject": "Weekly Brief",
		"editorial": "Hello world",
		"highlights": [
			{
				"title": "Highlight",
				"source_name": "Source",
				"source_url": "https://example.com",
				"analysis": "Analysis"
			}
		],
		"resources": []
	}`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	nl, err := LoadJSON(path)
	if err != nil {
		t.Fatalf("LoadJSON() error: %v", err)
	}

	if nl.Subject != "Weekly Brief" {
		t.Errorf("subject = %q, want %q", nl.Subject, "Weekly Brief")
	}
	if len(nl.Highlights) != 1 {
		t.Errorf("highlights length = %d, want 1", len(nl.Highlights))
	}
}

func TestLoadJSONInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")

	if err := os.WriteFile(path, []byte(`{"subject":`), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadJSON(path); err == nil {
		t.Fatal("expected LoadJSON() to fail on invalid JSON")
	}
}

func TestGetLatestOutputPath(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "2026-03-17.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "2026-03-18.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := GetLatestOutputPath(dir)
	if err != nil {
		t.Fatalf("GetLatestOutputPath() error: %v", err)
	}

	want := filepath.Join(dir, "2026-03-18.json")
	if got != want {
		t.Errorf("latest path = %q, want %q", got, want)
	}
}

func TestGetLatestOutputPathEmptyDir(t *testing.T) {
	dir := t.TempDir()

	if _, err := GetLatestOutputPath(dir); err == nil {
		t.Fatal("expected GetLatestOutputPath() to fail on empty output dir")
	}
}