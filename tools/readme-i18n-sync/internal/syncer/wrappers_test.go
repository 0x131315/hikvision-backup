package syncer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTranslateMissingWrapperUsesDefaultTranslator(t *testing.T) {
	prev := defaultTranslator
	stub := &stubTranslator{prefix: "ru:"}
	defaultTranslator = stub
	t.Cleanup(func() { defaultTranslator = prev })

	got, err := translateMissing(language{Code: "ru"}, []string{"a"}, false)
	if err != nil {
		t.Fatalf("translateMissing: %v", err)
	}
	if len(got) != 1 || got[0] != "ru:a" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestEnvTranslatorInitMode(t *testing.T) {
	tr := envTranslator{}
	got, err := tr.Translate(language{Code: "ru"}, []string{"a", "b"}, true)
	if err != nil {
		t.Fatalf("Translate init mode: %v", err)
	}
	if strings.Join(got, ",") != "a,b" {
		t.Fatalf("unexpected init-mode output: %v", got)
	}
}

func TestEnvTranslatorNoProvider(t *testing.T) {
	t.Setenv("DEEPL_API_KEY", "")
	t.Setenv("GOOGLE_TRANSLATE_API_KEY", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("LIBRETRANSLATE_URL", "")

	_, err := envTranslator{}.Translate(language{Code: "ru"}, []string{"a"}, false)
	if err == nil {
		t.Fatal("expected missing provider error")
	}
}

func TestProcessLanguageWrapperAndRun(t *testing.T) {
	tmp := t.TempDir()
	prevDefault := defaultTranslator
	prevI18n := i18nDir
	prevTM := tmDir
	prevSource := sourceFile
	prevLangs := langs

	defaultTranslator = &stubTranslator{prefix: "ru:"}
	i18nDir = filepath.Join(tmp, "i18n")
	tmDir = filepath.Join(i18nDir, "tm")
	sourceFile = filepath.Join(tmp, "README.md")
	langs = []language{{Code: "ru", TargetLang: "RU", GoogleLang: "ru"}}
	t.Cleanup(func() {
		defaultTranslator = prevDefault
		i18nDir = prevI18n
		tmDir = prevTM
		sourceFile = prevSource
		langs = prevLangs
	})

	if err := os.WriteFile(sourceFile, []byte("Languages: English\n\nHello"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := run(false, false, false); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	out, err := os.ReadFile(filepath.Join(i18nDir, "README.ru.md"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(out), "ru:Hello") {
		t.Fatalf("missing translated text: %q", string(out))
	}
}

func TestFillWrappers(t *testing.T) {
	prev := defaultTranslator
	stub := &stubTranslator{prefix: "ru:"}
	defaultTranslator = stub
	t.Cleanup(func() { defaultTranslator = prev })

	tm := tmFile{Blocks: map[string]string{}}
	table := "| K | D |\n|---|---|\n| A | Hello world |"
	if err := fillTableTranslations(language{Code: "ru"}, []string{table}, tm, false, false); err != nil {
		t.Fatalf("fillTableTranslations: %v", err)
	}
	if len(tm.Blocks) != 1 {
		t.Fatalf("expected table block in tm: %v", tm.Blocks)
	}

	blocks := []string{"Hello"}
	if err := fillMissingTranslations(language{Code: "ru"}, blocks, []int{0}, blocks, tm, false, false); err != nil {
		t.Fatalf("fillMissingTranslations: %v", err)
	}
	if tm.Blocks[hashString("Hello")] != "ru:Hello" {
		t.Fatalf("missing translated block: %v", tm.Blocks)
	}
}
