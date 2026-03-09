package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
)

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

func callGoogleTranslateWithServiceAccount(endpoint, credsPath string, texts []string, targetLang string) ([]string, error) {
	endpoint = strings.TrimRight(endpoint, "/")
	data, err := os.ReadFile(credsPath)
	if err != nil {
		return nil, fmt.Errorf("read google credentials: %w", err)
	}

	creds, err := google.CredentialsFromJSON(nil, data, "https://www.googleapis.com/auth/cloud-translation")
	if err != nil {
		return nil, fmt.Errorf("parse google credentials: %w", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("google token: %w", err)
	}

	form := url.Values{}
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
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

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
