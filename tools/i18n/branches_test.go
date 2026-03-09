package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunReadSourceError(t *testing.T) {
	prevSource := sourceFile
	sourceFile = filepath.Join(t.TempDir(), "missing.md")
	t.Cleanup(func() { sourceFile = prevSource })

	if err := run(false, false, false); err == nil {
		t.Fatal("expected read source error")
	}
}

func TestFillMissingTranslationsQuotaAndGenericError(t *testing.T) {
	blocks := []string{"A"}
	idx := []int{0}
	texts := []string{"A"}
	tm := tmFile{Blocks: map[string]string{}}

	err := fillMissingTranslationsWithTranslator(&stubTranslator{err: errors.New("quota exceeded")}, language{Code: "ru"}, blocks, idx, texts, tm, false, false)
	if err == nil || !strings.Contains(err.Error(), "translation quota exceeded") {
		t.Fatalf("expected quota-friendly error, got: %v", err)
	}

	err = fillMissingTranslationsWithTranslator(&stubTranslator{err: errors.New("boom")}, language{Code: "ru"}, blocks, idx, texts, tm, false, false)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected passthrough error, got: %v", err)
	}
}

func TestMarkdownFalseBranches(t *testing.T) {
	if isMarkdownTableBlock("plain text") {
		t.Fatal("plain text must not be a markdown table")
	}
	if got := firstNonEmptyLine([]string{"", "  "}, 0); got != "" {
		t.Fatalf("expected empty line result, got %q", got)
	}
	if got := findLanguagesBlock([]string{"Hello", "World"}); got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
	if isNumericLike("12a") {
		t.Fatal("alphanumeric must not be numeric-like")
	}
}

func TestWriteTMError(t *testing.T) {
	base := t.TempDir()
	parentAsFile := filepath.Join(base, "notdir")
	if err := os.WriteFile(parentAsFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	badPath := filepath.Join(parentAsFile, "tm.json")
	if err := writeTM(badPath, tmFile{Blocks: map[string]string{}}); err == nil {
		t.Fatal("expected writeTM mkdir error")
	}
}
