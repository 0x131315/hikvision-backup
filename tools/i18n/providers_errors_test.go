package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCallDeepLErrorPaths(t *testing.T) {
	quota := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(456)
		fmt.Fprint(w, "quota")
	}))
	defer quota.Close()
	if _, err := callDeepL(quota.URL, "k", []string{"a"}, "RU"); err == nil {
		t.Fatal("expected quota error")
	}

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad")
	}))
	defer bad.Close()
	if _, err := callDeepL(bad.URL, "k", []string{"a"}, "RU"); err == nil {
		t.Fatal("expected deepl non-2xx error")
	}
}

func TestCallGoogleTranslateErrorPaths(t *testing.T) {
	non2xx := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad")
	}))
	defer non2xx.Close()
	if _, err := callGoogleTranslate(non2xx.URL, "k", []string{"a"}, "ru"); err == nil {
		t.Fatal("expected google non-2xx error")
	}

	countMismatch := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":{"translations":[]}}`)
	}))
	defer countMismatch.Close()
	if _, err := callGoogleTranslate(countMismatch.URL, "k", []string{"a"}, "ru"); err == nil {
		t.Fatal("expected translation count mismatch error")
	}
}

func TestCallLibreTranslateErrorPath(t *testing.T) {
	non2xx := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad")
	}))
	defer non2xx.Close()
	if _, err := callLibreTranslate(non2xx.URL, "", []string{"a"}, "ru"); err == nil {
		t.Fatal("expected libre non-2xx error")
	}
}

func TestGoogleServiceAccountParseError(t *testing.T) {
	credsPath := filepath.Join(t.TempDir(), "creds.json")
	if err := os.WriteFile(credsPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("write creds: %v", err)
	}

	_, err := callGoogleTranslateWithServiceAccount("http://127.0.0.1/unused", credsPath, []string{"a"}, "ru")
	if err == nil {
		t.Fatal("expected parse google credentials error")
	}
}
