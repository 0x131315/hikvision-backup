package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func processLanguage(lang language, blocks, seps []string, sourceHash string, checkOnly, initMode, force bool) error {
	tmPath := filepath.Join(tmDir, fmt.Sprintf("README.%s.json", lang.Code))
	outPath := filepath.Join(i18nDir, fmt.Sprintf("README.%s.md", lang.Code))

	tm, err := loadTM(tmPath)
	if err != nil {
		return fmt.Errorf("load TM %s: %w", tmPath, err)
	}

	if tm.Blocks == nil {
		tm.Blocks = make(map[string]string)
	}

	// Positional sync is only safe for bootstrap (empty TM).
	// When source changes, syncing by index can mask missing translations.
	if tm.SourceHash == "" && len(tm.Blocks) == 0 {
		if err := syncFromTranslation(outPath, blocks, tm); err != nil {
			return err
		}
	}

	if force {
		tm.Blocks = make(map[string]string)
	}

	if idx := findLanguagesBlock(blocks); idx >= 0 {
		tm.Blocks[hashString(blocks[idx])] = languageSwitcherBlock(lang.Code)
	}

	if err := fillTableTranslations(lang, blocks, tm, initMode, force); err != nil {
		return err
	}

	missingIdx, missingText := findMissing(blocks, tm)
	if err := fillMissingTranslations(lang, blocks, missingIdx, missingText, tm, checkOnly, initMode); err != nil {
		return err
	}

	if checkOnly {
		return nil
	}

	tm.SourceHash = sourceHash
	tm.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := writeTM(tmPath, tm); err != nil {
		return err
	}

	out := buildTranslated(blocks, seps, tm)
	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}

	return nil
}

func fillTableTranslations(lang language, blocks []string, tm tmFile, initMode, force bool) error {
	for _, block := range blocks {
		if !isMarkdownTableBlock(block) {
			continue
		}

		key := hashString(block)
		if !force && tm.Blocks[key] != "" {
			continue
		}

		translatedTable, err := translateTableBlock(lang, block, initMode)
		if err != nil {
			return err
		}
		tm.Blocks[key] = translatedTable
	}

	return nil
}

func fillMissingTranslations(lang language, blocks []string, missingIdx []int, missingText []string, tm tmFile, checkOnly, initMode bool) error {
	if len(missingIdx) == 0 {
		return nil
	}
	if checkOnly {
		return fmt.Errorf("missing translations for %s: %d block(s)", lang.Code, len(missingIdx))
	}

	translated, err := translateMissing(lang, missingText, initMode)
	if err != nil {
		if isQuotaExceeded(err) {
			if len(missingIdx) > 0 {
				return fmt.Errorf("translation quota exceeded for %s and %d block(s) are missing; update manually", lang.Code, len(missingIdx))
			}
			return nil
		}
		return err
	}

	for i, idx := range missingIdx {
		tm.Blocks[hashString(blocks[idx])] = translated[i]
	}

	return nil
}

func findMissing(blocks []string, tm tmFile) ([]int, []string) {
	var idxs []int
	var texts []string
	for i, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		key := hashString(block)
		if _, ok := tm.Blocks[key]; ok {
			continue
		}
		idxs = append(idxs, i)
		texts = append(texts, block)
	}
	return idxs, texts
}

func buildTranslated(blocks, seps []string, tm tmFile) string {
	var buf bytes.Buffer
	for i, block := range blocks {
		if strings.TrimSpace(block) == "" {
			buf.WriteString(block)
		} else {
			key := hashString(block)
			if translated, ok := tm.Blocks[key]; ok {
				buf.WriteString(translated)
			} else {
				buf.WriteString(block)
			}
		}

		if i < len(seps) {
			buf.WriteString(seps[i])
		}
	}
	return buf.String()
}

func syncFromTranslation(path string, blocks []string, tm tmFile) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}

	transBlocks, _ := splitBlocks(string(data))
	if len(transBlocks) != len(blocks) {
		return nil
	}

	for i, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		translated := transBlocks[i]
		if strings.TrimSpace(translated) == "" {
			continue
		}
		tm.Blocks[hashString(block)] = translated
	}

	return nil
}
