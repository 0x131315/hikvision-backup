package syncer

import (
	"strings"
	"testing"
)

func TestLanguageSwitcherBlock(t *testing.T) {
	en := languageSwitcherBlock("en")
	if !strings.Contains(en, "Languages: English") {
		t.Fatalf("unexpected en switcher: %q", en)
	}
	if !strings.Contains(en, "(i18n/README.ru.md)") {
		t.Fatalf("expected ru link from en: %q", en)
	}

	ru := languageSwitcherBlock("ru")
	if !strings.Contains(ru, "../README.md") {
		t.Fatalf("expected en link from ru: %q", ru)
	}
	if !strings.Contains(ru, "README.zh.md") {
		t.Fatalf("expected zh link from ru: %q", ru)
	}
}

func TestSplitBlocks(t *testing.T) {
	blocks, seps := splitBlocks("a\n\n\nb\n\n\nc")
	if len(blocks) != 3 || len(seps) != 2 {
		t.Fatalf("unexpected split result: blocks=%d seps=%d", len(blocks), len(seps))
	}
	if blocks[1] != "b" {
		t.Fatalf("unexpected middle block: %q", blocks[1])
	}
}

func TestParseAndBuildMarkdownTable(t *testing.T) {
	input := "| A | B |\n|---|---|\n| x | y |"
	h, rows, ok := parseMarkdownTable(input)
	if !ok {
		t.Fatal("expected table to be parsed")
	}
	if len(h) != 2 || len(rows) != 1 || len(rows[0]) != 2 {
		t.Fatalf("unexpected table shape: h=%v rows=%v", h, rows)
	}

	out := buildMarkdownTable(h, rows)
	if !strings.Contains(out, "| A | B |") || !strings.Contains(out, "| x | y |") {
		t.Fatalf("unexpected built table: %q", out)
	}
}

func TestShouldTranslateTableCell(t *testing.T) {
	cases := []struct {
		cell string
		want bool
	}{
		{cell: "CAM_HOST", want: false},
		{cell: "--force", want: false},
		{cell: "true", want: false},
		{cell: "Host address", want: true},
		{cell: "`CAM_HOST`", want: false},
	}

	for _, tc := range cases {
		got := shouldTranslateTableCell(tc.cell)
		if got != tc.want {
			t.Fatalf("cell %q: got %v want %v", tc.cell, got, tc.want)
		}
	}
}

func TestTranslateTableBlockWithTranslator(t *testing.T) {
	block := "| Key | Description |\n|---|---|\n| CAM_HOST | Host address |\n| true | Keep bool |"
	tr := &stubTranslator{prefix: "ru:"}

	out, err := translateTableBlockWithTranslator(tr, language{Code: "ru"}, block, false)
	if err != nil {
		t.Fatalf("translate table: %v", err)
	}

	if tr.calls != 1 {
		t.Fatalf("expected 1 translator call, got %d", tr.calls)
	}
	if strings.Contains(out, "ru:CAM_HOST") {
		t.Fatalf("technical token should not be translated: %q", out)
	}
	if !strings.Contains(out, "ru:Host address") {
		t.Fatalf("natural language cell should be translated: %q", out)
	}
	if !strings.Contains(out, "|---|---|") {
		t.Fatalf("table separator should be preserved: %q", out)
	}
}

func TestTranslateTableBlockWrapperUsesDefaultTranslator(t *testing.T) {
	prev := defaultTranslator
	defaultTranslator = &stubTranslator{prefix: "ru:"}
	t.Cleanup(func() { defaultTranslator = prev })

	block := "| Key | Description |\n|---|---|\n| A | Simple text |"
	out, err := translateTableBlock(language{Code: "ru"}, block, false)
	if err != nil {
		t.Fatalf("translateTableBlock: %v", err)
	}
	if !strings.Contains(out, "ru:Simple text") {
		t.Fatalf("expected translated cell in wrapper output: %q", out)
	}
}
