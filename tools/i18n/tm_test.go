package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTMNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	tm, err := loadTM(path)
	if err != nil {
		t.Fatalf("loadTM not-exist: %v", err)
	}
	if tm.SourceHash != "" || tm.Blocks != nil {
		t.Fatalf("unexpected tm: %#v", tm)
	}
}

func TestLoadTMInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write bad json: %v", err)
	}
	if _, err := loadTM(path); err == nil {
		t.Fatal("expected json parse error")
	}
}

func TestWriteAndLoadTM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tm", "README.ru.json")
	in := tmFile{SourceHash: "h", UpdatedAt: "now", Blocks: map[string]string{"k": "v"}}
	if err := writeTM(path, in); err != nil {
		t.Fatalf("writeTM: %v", err)
	}

	out, err := loadTM(path)
	if err != nil {
		t.Fatalf("loadTM: %v", err)
	}
	if out.SourceHash != "h" || out.UpdatedAt != "now" || out.Blocks["k"] != "v" {
		t.Fatalf("unexpected tm readback: %#v", out)
	}
}
