package sources

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testLookback = 48 * time.Hour

func TestGetLastExecutionDate(t *testing.T) {
	dir := t.TempDir()

	content := `{"generated_at": "2026-03-10T08:00:00Z", "subject": "Test"}`
	if err := os.WriteFile(filepath.Join(dir, "2026-03-10.json"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := GetLastExecutionDate(dir, testLookback)
	if err != nil {
		t.Fatalf("GetLastExecutionDate() error: %v", err)
	}

	expected := time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestGetLastExecutionDateMultipleFiles(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "2026-03-08.json"), []byte(`{"generated_at": "2026-03-08T08:00:00Z"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "2026-03-10.json"), []byte(`{"generated_at": "2026-03-10T10:00:00Z"}`), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := GetLastExecutionDate(dir, testLookback)
	if err != nil {
		t.Fatalf("GetLastExecutionDate() error: %v", err)
	}

	expected := time.Date(2026, 3, 10, 10, 0, 0, 0, time.UTC)
	if !got.Equal(expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func TestGetLastExecutionDateEmptyDir(t *testing.T) {
	dir := t.TempDir()

	got, err := GetLastExecutionDate(dir, testLookback)
	if err != nil {
		t.Fatalf("GetLastExecutionDate() error: %v", err)
	}

	// Should be approximately now minus 48h
	expected := time.Now().UTC().Add(-testLookback)
	diff := got.Sub(expected)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expected ~%v, got %v (diff: %v)", expected, got, diff)
	}
}
