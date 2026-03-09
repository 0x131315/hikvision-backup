package config

import (
	"log/slog"
	"testing"
)

func TestBuildBaseURL(t *testing.T) {
	got, err := buildBaseURL("192.168.1.10")
	if err != nil {
		t.Fatalf("buildBaseURL error: %v", err)
	}
	if got != "https://192.168.1.10" {
		t.Fatalf("unexpected base url: %s", got)
	}

	got, err = buildBaseURL("http://cam.local/")
	if err != nil {
		t.Fatalf("buildBaseURL error: %v", err)
	}
	if got != "http://cam.local" {
		t.Fatalf("unexpected base url trim: %s", got)
	}
}

func TestInitMissingRequiredEnv(t *testing.T) {
	t.Setenv("CAM_HOST", "")
	t.Setenv("CAM_USER", "admin")
	t.Setenv("DOWNLOAD_DIR", t.TempDir())

	_, err := Init(slog.LevelInfo, false)
	if err == nil {
		t.Fatalf("expected missing CAM_HOST error")
	}
}

func TestInitParsesConfig(t *testing.T) {
	t.Setenv("CAM_HOST", "192.168.1.10")
	t.Setenv("CAM_USER", "admin")
	t.Setenv("CAM_PASS", "secret")
	t.Setenv("CAM_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("NO_PROXY", "true")
	t.Setenv("DOWNLOAD_DIR", t.TempDir())
	t.Setenv("SCAN_LAST_DAYS", "-3")
	t.Setenv("SCAN_FROM_LOCAL_LATEST", "0")
	t.Setenv("HTTP_RETRY_CNT", "5")
	t.Setenv("HTTP_TIMEOUT", "10")

	conf, err := Init(slog.LevelDebug, true)
	if err != nil {
		t.Fatalf("Init error: %v", err)
	}

	if conf.BaseURL != "https://192.168.1.10" {
		t.Fatalf("unexpected BaseURL: %s", conf.BaseURL)
	}
	if !conf.NoProxy {
		t.Fatalf("expected NoProxy=true")
	}
	if !conf.InsecureTLS {
		t.Fatalf("expected InsecureTLS=true")
	}
	if conf.ScanLastDays != 3 {
		t.Fatalf("expected ScanLastDays=3, got %d", conf.ScanLastDays)
	}
	if conf.RetryCnt != 5 || conf.HttpTimeout != 10 {
		t.Fatalf("unexpected retry/timeout: retry=%d timeout=%d", conf.RetryCnt, conf.HttpTimeout)
	}
}
