package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type language struct {
	Code       string
	TargetLang string
	GoogleLang string
}

type tmFile struct {
	SourceHash string            `json:"source_hash"`
	UpdatedAt  string            `json:"updated_at"`
	Blocks     map[string]string `json:"blocks"`
}

var (
	sourceFile = "README.md"
	i18nDir    = "i18n"
	tmDir      = filepath.Join(i18nDir, "tm")
	langs      = []language{
		{Code: "ru", TargetLang: "RU", GoogleLang: "ru"},
		{Code: "zh", TargetLang: "ZH", GoogleLang: "zh-CN"},
	}
	sepRe = regexp.MustCompile(`\n{2,}`)
)

func main() {
	checkOnly := flag.Bool("check", false, "check translations without updating")
	initMode := flag.Bool("init", false, "bootstrap translations using source text when missing")
	force := flag.Bool("force", false, "retranslate all blocks and overwrite existing translations")
	flag.Parse()

	source, err := os.ReadFile(sourceFile)
	exitOnErr("read README.md", err)

	blocks, seps := splitBlocks(string(source))
	sourceHash := hashString(string(source))

	for _, lang := range langs {
		if err := processLanguage(lang, blocks, seps, sourceHash, *checkOnly, *initMode, *force); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

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

	if err := syncFromTranslation(outPath, blocks, seps, tm); err != nil {
		return err
	}

	if force {
		tm.Blocks = make(map[string]string)
	}

	missingIdx, missingText := findMissing(blocks, tm)

	if len(missingIdx) > 0 {
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

func syncFromTranslation(path string, blocks []string, seps []string, tm tmFile) error {
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

func translateMissing(lang language, texts []string, initMode bool) ([]string, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	if initMode {
		return texts, nil
	}

	if key := os.Getenv("DEEPL_API_KEY"); key != "" {
		endpoint := os.Getenv("DEEPL_API_URL")
		if endpoint == "" {
			endpoint = deeplEndpointForKey(key)
		}
		return callDeepL(endpoint, key, texts, lang.TargetLang)
	}

	if key := os.Getenv("GOOGLE_TRANSLATE_API_KEY"); key != "" {
		endpoint := os.Getenv("GOOGLE_TRANSLATE_API_URL")
		if endpoint == "" {
			endpoint = "https://translation.googleapis.com/language/translate/v2"
		}
		return callGoogleTranslate(endpoint, key, texts, lang.GoogleLang)
	}

	libreURL := os.Getenv("LIBRETRANSLATE_URL")
	if libreURL == "" {
		return nil, fmt.Errorf("DEEPL_API_KEY is not set, GOOGLE_TRANSLATE_API_KEY is not set, and LIBRETRANSLATE_URL is empty")
	}
	libreKey := os.Getenv("LIBRETRANSLATE_API_KEY")
	return callLibreTranslate(libreURL, libreKey, texts, lang.Code)
}

func callDeepL(endpoint, key string, texts []string, targetLang string) ([]string, error) {
	form := url.Values{}
	form.Set("auth_key", key)
	form.Set("target_lang", targetLang)
	for _, text := range texts {
		form.Add("text", text)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 456 || resp.StatusCode == 429 {
		return nil, fmt.Errorf("quota exceeded")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("deepl error: %s", strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Translations) != len(texts) {
		return nil, fmt.Errorf("unexpected translation count: %d", len(parsed.Translations))
	}

	out := make([]string, len(texts))
	for i, tr := range parsed.Translations {
		out[i] = tr.Text
	}
	return out, nil
}

func deeplEndpointForKey(key string) string {
	if strings.HasSuffix(key, ":fx") {
		return "https://api-free.deepl.com/v2/translate"
	}
	return "https://api.deepl.com/v2/translate"
}

func isQuotaExceeded(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "quota exceeded")
}

func callLibreTranslate(endpoint, key string, texts []string, targetLang string) ([]string, error) {
	endpoint = strings.TrimRight(endpoint, "/") + "/translate"

	out := make([]string, len(texts))
	for i, text := range texts {
		payload := map[string]string{
			"q":      text,
			"source": "en",
			"target": targetLang,
			"format": "text",
		}
		if key != "" {
			payload["api_key"] = key
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("libretranslate error: %s", strings.TrimSpace(string(body)))
		}

		var parsed struct {
			TranslatedText string `json:"translatedText"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, err
		}
		out[i] = parsed.TranslatedText
	}
	return out, nil
}

func callGoogleTranslate(endpoint, key string, texts []string, targetLang string) ([]string, error) {
	endpoint = strings.TrimRight(endpoint, "/")
	form := url.Values{}
	form.Set("key", key)
	form.Set("target", targetLang)
	form.Set("format", "text")
	for _, text := range texts {
		form.Add("q", text)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("google translate error: %s", strings.TrimSpace(string(body)))
	}

	var parsed struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Data.Translations) != len(texts) {
		return nil, fmt.Errorf("unexpected translation count: %d", len(parsed.Data.Translations))
	}

	out := make([]string, len(texts))
	for i, tr := range parsed.Data.Translations {
		out[i] = htmlUnescape(tr.TranslatedText)
	}
	return out, nil
}

func htmlUnescape(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
	)
	return replacer.Replace(s)
}

func loadTM(path string) (tmFile, error) {
	var tm tmFile
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return tm, nil
		}
		return tm, err
	}
	if err := json.Unmarshal(data, &tm); err != nil {
		return tm, err
	}
	return tm, nil
}

func writeTM(path string, tm tmFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum[:])
}

func exitOnErr(action string, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %v\n", action, err)
	os.Exit(1)
}
