package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "videos")

	if err := Init(target); err != nil {
		t.Fatalf("Init error: %v", err)
	}
	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("stat target dir: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected directory, got file")
	}
}

func TestFileSizeAndExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "video.mp4")
	data := []byte("12345")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	if !IsFileExist(path) {
		t.Fatalf("expected file to exist")
	}
	if got := FileSize(path); got != len(data) {
		t.Fatalf("expected size %d, got %d", len(data), got)
	}
}

func TestIsFileExistReturnsFalseForDirectory(t *testing.T) {
	dir := t.TempDir()
	if IsFileExist(dir) {
		t.Fatalf("directory must not be treated as file")
	}
}

func TestRemoveFileMissingDoesNotFail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.mp4")

	RemoveFile(path)
	if IsFileExist(path) {
		t.Fatalf("file should not exist")
	}
}
