package main

import (
	"fmt"
	"os"
)

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

	if credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credsPath != "" {
		endpoint := os.Getenv("GOOGLE_TRANSLATE_API_URL")
		if endpoint == "" {
			endpoint = "https://translation.googleapis.com/language/translate/v2"
		}
		return callGoogleTranslateWithServiceAccount(endpoint, credsPath, texts, lang.GoogleLang)
	}

	libreURL := os.Getenv("LIBRETRANSLATE_URL")
	if libreURL == "" {
		return nil, fmt.Errorf("DEEPL_API_KEY is not set, GOOGLE_TRANSLATE_API_KEY is not set, GOOGLE_APPLICATION_CREDENTIALS is not set, and LIBRETRANSLATE_URL is empty")
	}
	libreKey := os.Getenv("LIBRETRANSLATE_API_KEY")
	return callLibreTranslate(libreURL, libreKey, texts, lang.Code)
}
