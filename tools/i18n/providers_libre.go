package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
