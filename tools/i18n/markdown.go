package main

import (
	"fmt"
	"strings"
)

func splitBlocks(s string) ([]string, []string) {
	idxs := sepRe.FindAllStringIndex(s, -1)
	if len(idxs) == 0 {
		return []string{s}, nil
	}

	blocks := make([]string, 0, len(idxs)+1)
	seps := make([]string, 0, len(idxs))
	start := 0
	for _, idx := range idxs {
		blocks = append(blocks, s[start:idx[0]])
		seps = append(seps, s[idx[0]:idx[1]])
		start = idx[1]
	}
	blocks = append(blocks, s[start:])
	return blocks, seps
}

func findLanguagesBlock(blocks []string) int {
	for i, block := range blocks {
		if strings.HasPrefix(strings.TrimSpace(block), "Languages:") {
			return i
		}
	}
	return -1
}

func languageSwitcherBlock(langCode string) string {
	current := langCode
	if current != "en" && current != "ru" && current != "zh" {
		current = "en"
	}

	parts := make([]string, 0, len(languageNav))
	for _, item := range languageNav {
		if item.Code == current {
			parts = append(parts, item.Name)
			continue
		}

		link := languageDocLink(current, item.Code)
		parts = append(parts, fmt.Sprintf("[%s](%s)", item.Name, link))
	}

	return "Languages: " + strings.Join(parts, " | ")
}

func languageDocLink(currentCode, targetCode string) string {
	if currentCode == "en" {
		return fmt.Sprintf("i18n/README.%s.md", targetCode)
	}
	if targetCode == "en" {
		return "../README.md"
	}
	return fmt.Sprintf("README.%s.md", targetCode)
}

func isMarkdownTableBlock(block string) bool {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) < 2 {
		return false
	}

	first := firstNonEmptyLine(lines, 0)
	second := firstNonEmptyLine(lines, 1)
	if first == "" || second == "" {
		return false
	}
	if !strings.Contains(first, "|") {
		return false
	}
	if !strings.Contains(second, "|") || !strings.Contains(second, "-") {
		return false
	}

	return true
}

func firstNonEmptyLine(lines []string, skip int) string {
	seen := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if seen < skip {
			seen++
			continue
		}
		return trimmed
	}
	return ""
}

func translateTableBlock(lang language, block string, initMode bool) (string, error) {
	return translateTableBlockWithTranslator(defaultTranslator, lang, block, initMode)
}

func translateTableBlockWithTranslator(translator Translator, lang language, block string, initMode bool) (string, error) {
	header, rows, ok := parseMarkdownTable(block)
	if !ok {
		return block, nil
	}

	texts := make([]string, 0, len(header)+len(rows)*len(header))
	targets := make([]*string, 0, len(header)+len(rows)*len(header))
	for i := range header {
		if shouldTranslateTableCell(header[i]) {
			texts = append(texts, header[i])
			targets = append(targets, &header[i])
		}
	}
	for ri := range rows {
		for ci := range rows[ri] {
			if shouldTranslateTableCell(rows[ri][ci]) {
				texts = append(texts, rows[ri][ci])
				targets = append(targets, &rows[ri][ci])
			}
		}
	}
	if len(texts) == 0 {
		return buildMarkdownTable(header, rows), nil
	}

	translated, err := translator.Translate(lang, texts, initMode)
	if err != nil {
		return "", err
	}
	if len(translated) != len(targets) {
		return "", fmt.Errorf("unexpected translated table cell count: %d", len(translated))
	}
	for i := range translated {
		*targets[i] = translated[i]
	}

	return buildMarkdownTable(header, rows), nil
}

func parseMarkdownTable(block string) ([]string, [][]string, bool) {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	if len(lines) < 2 {
		return nil, nil, false
	}

	headerIdx := -1
	sepIdx := -1
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.Contains(line, "|") {
			continue
		}
		if headerIdx == -1 {
			headerIdx = i
			continue
		}
		if strings.Contains(line, "-") {
			sepIdx = i
			break
		}
	}
	if headerIdx == -1 || sepIdx == -1 || sepIdx <= headerIdx {
		return nil, nil, false
	}

	header := splitTableCells(lines[headerIdx])
	if len(header) == 0 {
		return nil, nil, false
	}

	rows := make([][]string, 0)
	for i := sepIdx + 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.Contains(line, "|") {
			continue
		}
		row := splitTableCells(line)
		if len(row) == 0 {
			continue
		}

		// Keep column count stable.
		if len(row) < len(header) {
			row = append(row, make([]string, len(header)-len(row))...)
		}
		if len(row) > len(header) {
			row = row[:len(header)]
		}
		rows = append(rows, row)
	}

	return header, rows, true
}

func splitTableCells(line string) []string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "|") {
		line = strings.TrimPrefix(line, "|")
	}
	if strings.HasSuffix(line, "|") {
		line = strings.TrimSuffix(line, "|")
	}

	parts := strings.Split(line, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func buildMarkdownTable(header []string, rows [][]string) string {
	var b strings.Builder
	b.WriteString(joinTableCells(header))
	b.WriteByte('\n')
	b.WriteString("|")
	for range header {
		b.WriteString("---|")
	}

	for _, row := range rows {
		b.WriteByte('\n')
		b.WriteString(joinTableCells(row))
	}

	return b.String()
}

func joinTableCells(cells []string) string {
	var b strings.Builder
	b.WriteString("|")
	for _, c := range cells {
		b.WriteString(" ")
		b.WriteString(strings.TrimSpace(c))
		b.WriteString(" |")
	}
	return b.String()
}

func shouldTranslateTableCell(cell string) bool {
	s := strings.TrimSpace(cell)
	if s == "" || s == "-" {
		return false
	}

	l := strings.ToLower(strings.Trim(s, " ."))
	switch l {
	case "true", "false", "yes", "no", "empty", "secret", "info", "debug":
		return false
	}
	if isNumericLike(l) {
		return false
	}

	// Single-token cells are often technical values and should stay as-is.
	// Multi-token cells are usually natural language descriptions and should be translated.
	fields := strings.Fields(s)
	if len(fields) == 1 && isLikelyTechnicalToken(fields[0]) {
		return false
	}

	return true
}

func isNumericLike(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && r != '.' {
			return false
		}
	}
	return true
}

func isLikelyTechnicalToken(s string) bool {
	token := strings.TrimSpace(s)
	if token == "" {
		return false
	}
	if strings.HasPrefix(token, "`") && strings.HasSuffix(token, "`") {
		return true
	}
	if strings.HasPrefix(token, "http://") || strings.HasPrefix(token, "https://") {
		return true
	}
	if strings.HasPrefix(token, "--") || (strings.HasPrefix(token, "-") && len(token) <= 4) {
		return true
	}
	if strings.Contains(token, "/") || strings.Contains(token, "_") {
		return true
	}
	return false
}
