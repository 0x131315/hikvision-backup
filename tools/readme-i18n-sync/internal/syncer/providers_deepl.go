package syncer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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
