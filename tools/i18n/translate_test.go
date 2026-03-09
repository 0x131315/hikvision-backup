package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEnvTranslatorDeepLBranch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"translations":[{"text":"ru:a"}]}`)
	}))
	defer ts.Close()

	t.Setenv("DEEPL_API_KEY", "k")
	t.Setenv("DEEPL_API_URL", ts.URL)
	t.Setenv("GOOGLE_TRANSLATE_API_KEY", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("LIBRETRANSLATE_URL", "")

	got, err := (envTranslator{}).Translate(language{Code: "ru", TargetLang: "RU"}, []string{"a"}, false)
	if err != nil {
		t.Fatalf("Translate deepl: %v", err)
	}
	if len(got) != 1 || got[0] != "ru:a" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestEnvTranslatorGoogleAPIKeyBranch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":{"translations":[{"translatedText":"ru:b"}]}}`)
	}))
	defer ts.Close()

	t.Setenv("DEEPL_API_KEY", "")
	t.Setenv("GOOGLE_TRANSLATE_API_KEY", "k")
	t.Setenv("GOOGLE_TRANSLATE_API_URL", ts.URL)
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("LIBRETRANSLATE_URL", "")

	got, err := (envTranslator{}).Translate(language{Code: "ru", GoogleLang: "ru"}, []string{"b"}, false)
	if err != nil {
		t.Fatalf("Translate google api key: %v", err)
	}
	if len(got) != 1 || got[0] != "ru:b" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestEnvTranslatorServiceAccountBranchReadError(t *testing.T) {
	t.Setenv("DEEPL_API_KEY", "")
	t.Setenv("GOOGLE_TRANSLATE_API_KEY", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/no/such/credentials.json")
	t.Setenv("GOOGLE_TRANSLATE_API_URL", "http://127.0.0.1/unused")
	t.Setenv("LIBRETRANSLATE_URL", "")

	_, err := (envTranslator{}).Translate(language{Code: "ru", GoogleLang: "ru"}, []string{"b"}, false)
	if err == nil {
		t.Fatal("expected error for missing service account file")
	}
}

func TestEnvTranslatorLibreBranch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"translatedText":"ru:c"}`)
	}))
	defer ts.Close()

	t.Setenv("DEEPL_API_KEY", "")
	t.Setenv("GOOGLE_TRANSLATE_API_KEY", "")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
	t.Setenv("LIBRETRANSLATE_URL", ts.URL)
	t.Setenv("LIBRETRANSLATE_API_KEY", "")

	got, err := (envTranslator{}).Translate(language{Code: "ru"}, []string{"c"}, false)
	if err != nil {
		t.Fatalf("Translate libre: %v", err)
	}
	if len(got) != 1 || got[0] != "ru:c" {
		t.Fatalf("unexpected result: %v", got)
	}
}
