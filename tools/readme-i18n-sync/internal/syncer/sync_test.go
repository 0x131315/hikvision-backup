package syncer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindMissing(t *testing.T) {
	blocks := []string{"", "A", "B"}
	tm := tmFile{Blocks: map[string]string{hashString("A"): "A-ru"}}

	idx, texts := findMissing(blocks, tm)
	if len(idx) != 1 || idx[0] != 2 {
		t.Fatalf("unexpected missing idx: %v", idx)
	}
	if len(texts) != 1 || texts[0] != "B" {
		t.Fatalf("unexpected missing texts: %v", texts)
	}
}

func TestBuildTranslated(t *testing.T) {
	blocks := []string{"A", "B"}
	seps := []string{"\n\n"}
	tm := tmFile{Blocks: map[string]string{hashString("A"): "A-ru"}}

	out := buildTranslated(blocks, seps, tm)
	if out != "A-ru\n\nB" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestProcessLanguageWithTranslator(t *testing.T) {
	tmp := t.TempDir()
	prevI18n := i18nDir
	prevTM := tmDir
	i18nDir = filepath.Join(tmp, "i18n")
	tmDir = filepath.Join(i18nDir, "tm")
	t.Cleanup(func() {
		i18nDir = prevI18n
		tmDir = prevTM
	})

	blocks := []string{
		"Languages: English | [Русский](i18n/README.ru.md) | [中文](i18n/README.zh.md)",
		"Hello world",
	}
	seps := []string{"\n\n"}
	tr := &stubTranslator{prefix: "ru:"}

	err := processLanguageWithTranslator(tr, language{Code: "ru"}, blocks, seps, hashString(strings.Join(blocks, "\n\n")), false, false, false)
	if err != nil {
		t.Fatalf("process language: %v", err)
	}

	outPath := filepath.Join(i18nDir, "README.ru.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "Languages:") {
		t.Fatalf("missing language switcher in output: %q", out)
	}
	if !strings.Contains(out, "ru:Hello world") {
		t.Fatalf("missing translated content: %q", out)
	}

	tmPath := filepath.Join(tmDir, "README.ru.json")
	if _, err := os.Stat(tmPath); err != nil {
		t.Fatalf("tm file was not written: %v", err)
	}
}

func TestFillMissingTranslationsCheckOnly(t *testing.T) {
	tm := tmFile{Blocks: map[string]string{}}
	blocks := []string{"A"}
	idx := []int{0}
	texts := []string{"A"}

	err := fillMissingTranslationsWithTranslator(&stubTranslator{prefix: "ru:"}, language{Code: "ru"}, blocks, idx, texts, tm, true, false)
	if err == nil {
		t.Fatal("expected check-only error")
	}
	if !strings.Contains(err.Error(), "missing translations") {
		t.Fatalf("unexpected error: %v", err)
	}
}
