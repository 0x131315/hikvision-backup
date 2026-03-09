package syncer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCallGoogleTranslateWithServiceAccountSuccess(t *testing.T) {
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		fmt.Fprint(w, `{"access_token":"tok","token_type":"Bearer","expires_in":3600}`)
	}))
	defer tokenSrv.Close()

	apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error":"bad auth: %s"}`, got)
			return
		}
		fmt.Fprint(w, `{"data":{"translations":[{"translatedText":"a &amp; b"}]}}`)
	}))
	defer apiSrv.Close()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8})

	credsJSON := fmt.Sprintf(`{
		"type":"service_account",
		"project_id":"test-project",
		"private_key_id":"kid",
		"private_key":%q,
		"client_email":"svc@test-project.iam.gserviceaccount.com",
		"token_uri":%q
	}`, string(pemKey), tokenSrv.URL)

	credsPath := filepath.Join(t.TempDir(), "creds.json")
	if err := os.WriteFile(credsPath, []byte(credsJSON), 0o644); err != nil {
		t.Fatalf("write creds: %v", err)
	}

	got, err := callGoogleTranslateWithServiceAccount(apiSrv.URL, credsPath, []string{"a"}, "ru")
	if err != nil {
		t.Fatalf("callGoogleTranslateWithServiceAccount: %v", err)
	}
	if len(got) != 1 || got[0] != "a & b" {
		t.Fatalf("unexpected translations: %v", got)
	}
}

func TestCallGoogleTranslateWithServiceAccountAPIError(t *testing.T) {
	tokenSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"access_token":"tok","token_type":"Bearer","expires_in":3600}`)
	}))
	defer tokenSrv.Close()

	apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"bad request"}`)
	}))
	defer apiSrv.Close()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8})
	credsJSON := fmt.Sprintf(`{"type":"service_account","private_key":%q,"client_email":"svc@test.iam.gserviceaccount.com","token_uri":%q}`, string(pemKey), tokenSrv.URL)

	credsPath := filepath.Join(t.TempDir(), "creds.json")
	if err := os.WriteFile(credsPath, []byte(credsJSON), 0o644); err != nil {
		t.Fatalf("write creds: %v", err)
	}

	_, err = callGoogleTranslateWithServiceAccount(apiSrv.URL, credsPath, []string{"a"}, "ru")
	if err == nil || !strings.Contains(err.Error(), "google translate error") {
		t.Fatalf("expected google translate error, got: %v", err)
	}
}
