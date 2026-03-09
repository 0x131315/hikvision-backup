package syncer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncFromTranslationCases(t *testing.T) {
	blocks := []string{"A", "B"}
	tm := tmFile{Blocks: map[string]string{}}

	if err := syncFromTranslation(filepath.Join(t.TempDir(), "missing.md"), blocks, tm); err != nil {
		t.Fatalf("syncFromTranslation missing file: %v", err)
	}

	path := filepath.Join(t.TempDir(), "README.ru.md")
	if err := os.WriteFile(path, []byte("only one block"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := syncFromTranslation(path, blocks, tm); err != nil {
		t.Fatalf("syncFromTranslation mismatch length: %v", err)
	}
	if len(tm.Blocks) != 0 {
		t.Fatalf("tm should stay empty on mismatch: %v", tm.Blocks)
	}

	if err := os.WriteFile(path, []byte("A-tr\n\nB-tr"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := syncFromTranslation(path, blocks, tm); err != nil {
		t.Fatalf("syncFromTranslation success: %v", err)
	}
	if tm.Blocks[hashString("A")] != "A-tr" || tm.Blocks[hashString("B")] != "B-tr" {
		t.Fatalf("unexpected tm blocks: %v", tm.Blocks)
	}
}

func TestEqualStringMaps(t *testing.T) {
	if !equalStringMaps(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Fatal("expected maps to be equal")
	}
	if equalStringMaps(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Fatal("expected maps to differ by value")
	}
	if equalStringMaps(map[string]string{"a": "1"}, map[string]string{"a": "1", "b": "2"}) {
		t.Fatal("expected maps to differ by length")
	}
}
