package config

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindLatestLocalFileTime(t *testing.T) {
	dir := t.TempDir()

	// Valid files
	t1 := time.Date(2024, 10, 1, 10, 0, 0, 0, time.Local)
	t2 := time.Date(2024, 10, 3, 12, 30, 0, 0, time.Local)
	if err := writeEmptyFile(filepath.Join(dir, t1.Format(localFileTimeLayout)+".mp4")); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := writeEmptyFile(filepath.Join(dir, t2.Format(localFileTimeLayout)+".mp4")); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// Invalid or irrelevant files
	_ = writeEmptyFile(filepath.Join(dir, "not-a-date.mp4"))
	_ = writeEmptyFile(filepath.Join(dir, "2024-10-04T10-00-00.txt"))

	got, ok := findLatestLocalFileTime(dir)
	if !ok {
		t.Fatalf("expected to find at least one valid file")
	}
	if !got.Equal(t2) {
		t.Fatalf("expected latest time %v, got %v", t2, got)
	}
}

func TestAdjustScanLastDaysFromLocal_NoFiles(t *testing.T) {
	dir := t.TempDir()

	got := adjustScanLastDaysFromLocal(dir, 10, 2)
	if got != 10 {
		t.Fatalf("expected scanLastDays unchanged, got %d", got)
	}
}

func TestAdjustScanLastDaysFromLocal_UsesCalculatedWhenEnvZero(t *testing.T) {
	dir := t.TempDir()

	latest := time.Now().AddDate(0, 0, -5).In(time.Local)
	if err := writeEmptyFile(filepath.Join(dir, latest.Format(localFileTimeLayout)+".mp4")); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// latest - 2 days
	target := latest.AddDate(0, 0, -2)
	expected := int(math.Ceil(time.Since(target).Hours() / 24))
	if expected < 1 {
		expected = 1
	}

	got := adjustScanLastDaysFromLocal(dir, 0, 2)
	if got != expected {
		t.Fatalf("expected %d, got %d", expected, got)
	}
}

func TestAdjustScanLastDaysFromLocal_TakesMinWithEnvValue(t *testing.T) {
	dir := t.TempDir()

	latest := time.Now().AddDate(0, 0, -10).In(time.Local)
	if err := writeEmptyFile(filepath.Join(dir, latest.Format(localFileTimeLayout)+".mp4")); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// latest - 2 days
	target := latest.AddDate(0, 0, -2)
	calculated := int(math.Ceil(time.Since(target).Hours() / 24))
	if calculated < 1 {
		calculated = 1
	}

	// Env value smaller than calculated should win
	if calculated <= 3 {
		t.Fatalf("calculated days too small for this test: %d", calculated)
	}

	got := adjustScanLastDaysFromLocal(dir, 3, 2)
	if got != 3 {
		t.Fatalf("expected min value 3, got %d", got)
	}
}

func TestAdjustScanLastDaysFromLocal_UsesBackfillDays(t *testing.T) {
	dir := t.TempDir()

	latest := time.Now().AddDate(0, 0, -10).In(time.Local)
	if err := writeEmptyFile(filepath.Join(dir, latest.Format(localFileTimeLayout)+".mp4")); err != nil {
		t.Fatalf("write file: %v", err)
	}

	backfill := 5
	target := latest.AddDate(0, 0, -backfill)
	expected := int(math.Ceil(time.Since(target).Hours() / 24))
	if expected < 1 {
		expected = 1
	}

	got := adjustScanLastDaysFromLocal(dir, 0, backfill)
	if got != expected {
		t.Fatalf("expected %d, got %d", expected, got)
	}
}

func writeEmptyFile(path string) error {
	return os.WriteFile(path, []byte{}, 0o644)
}
