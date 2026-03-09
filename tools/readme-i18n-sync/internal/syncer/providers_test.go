package syncer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeepLEndpointForKey(t *testing.T) {
	if got := deeplEndpointForKey("abc:fx"); got != "https://api-free.deepl.com/v2/translate" {
		t.Fatalf("unexpected free endpoint: %s", got)
	}
	if got := deeplEndpointForKey("abc"); got != "https://api.deepl.com/v2/translate" {
		t.Fatalf("unexpected paid endpoint: %s", got)
	}
}

func TestCallDeepL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("auth_key") != "k" || r.Form.Get("target_lang") != "RU" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"bad request"}`)
			return
		}
		fmt.Fprint(w, `{"translations":[{"text":"ru:one"},{"text":"ru:two"}]}`)
	}))
	defer ts.Close()

	got, err := callDeepL(ts.URL, "k", []string{"one", "two"}, "RU")
	if err != nil {
		t.Fatalf("callDeepL failed: %v", err)
	}
	if len(got) != 2 || got[0] != "ru:one" || got[1] != "ru:two" {
		t.Fatalf("unexpected translations: %v", got)
	}
}

func TestCallGoogleTranslate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"data":{"translations":[{"translatedText":"a &amp; b"}]}}`)
	}))
	defer ts.Close()

	got, err := callGoogleTranslate(ts.URL, "k", []string{"one"}, "ru")
	if err != nil {
		t.Fatalf("callGoogleTranslate failed: %v", err)
	}
	if len(got) != 1 || got[0] != "a & b" {
		t.Fatalf("unexpected translations: %v", got)
	}
}

func TestCallLibreTranslate(t *testing.T) {
	requests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		fmt.Fprint(w, `{"translatedText":"ru:text"}`)
	}))
	defer ts.Close()

	got, err := callLibreTranslate(ts.URL, "", []string{"one", "two"}, "ru")
	if err != nil {
		t.Fatalf("callLibreTranslate failed: %v", err)
	}
	if len(got) != 2 || got[0] != "ru:text" || got[1] != "ru:text" {
		t.Fatalf("unexpected translations: %v", got)
	}
	if requests != 2 {
		t.Fatalf("expected one request per text, got %d", requests)
	}
}
