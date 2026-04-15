package main

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/0x131315/hikvision-backup/internal/app/config"
	"github.com/spf13/cobra"
)

func TestApplyConfigFlagOverridesUsesCLIValues(t *testing.T) {
	downloadDir := setBaseConfigEnv(t)
	t.Setenv("CAM_HOST", "env-camera")
	t.Setenv("CAM_USER", "env-user")
	t.Setenv("CAM_PASS", "env-pass")
	t.Setenv("CAM_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("NO_PROXY", "true")
	t.Setenv("SCAN_LAST_DAYS", "9")
	t.Setenv("SCAN_FROM_LOCAL_LATEST", "4")
	t.Setenv("HTTP_RETRY_CNT", "5")
	t.Setenv("HTTP_TIMEOUT", "30")

	cmd, overrides := newConfigOverrideTestCommand()
	if err := cmd.ParseFlags([]string{
		"--download-dir=" + downloadDir,
		"--cam-host=cli-camera",
		"--cam-user=cli-user",
		"--cam-pass=cli-pass",
		"--cam-insecure-skip-verify=false",
		"--scan-last-days=0",
		"--scan-from-local-latest=0",
		"--http-retry-cnt=1",
		"--http-timeout=45",
		"--no-proxy=false",
	}); err != nil {
		t.Fatalf("ParseFlags error: %v", err)
	}

	if err := applyConfigFlagOverrides(cmd.Flags(), overrides); err != nil {
		t.Fatalf("applyConfigFlagOverrides error: %v", err)
	}

	conf, err := config.Init(slog.LevelInfo, false)
	if err != nil {
		t.Fatalf("Init error: %v", err)
	}

	if conf.Host != "cli-camera" {
		t.Fatalf("expected CLI CAM_HOST override, got %q", conf.Host)
	}
	if conf.User != "cli-user" {
		t.Fatalf("expected CLI CAM_USER override, got %q", conf.User)
	}
	if conf.Pass != "cli-pass" {
		t.Fatalf("expected CLI CAM_PASS override, got %q", conf.Pass)
	}
	if conf.InsecureTLS {
		t.Fatalf("expected CLI CAM_INSECURE_SKIP_VERIFY=false override")
	}
	if conf.NoProxy {
		t.Fatalf("expected CLI NO_PROXY=false override")
	}
	if conf.ScanLastDays != 0 {
		t.Fatalf("expected CLI SCAN_LAST_DAYS override, got %d", conf.ScanLastDays)
	}
	if conf.ScanLocal != 0 {
		t.Fatalf("expected CLI SCAN_FROM_LOCAL_LATEST override, got %d", conf.ScanLocal)
	}
	if conf.RetryCnt != 1 {
		t.Fatalf("expected CLI HTTP_RETRY_CNT override, got %d", conf.RetryCnt)
	}
	if conf.HttpTimeout != 45 {
		t.Fatalf("expected CLI HTTP_TIMEOUT override, got %d", conf.HttpTimeout)
	}
}

func TestApplyConfigFlagOverridesKeepsEnvWhenFlagsAreMissing(t *testing.T) {
	downloadDir := setBaseConfigEnv(t)
	t.Setenv("CAM_HOST", "env-camera")
	t.Setenv("CAM_USER", "env-user")
	t.Setenv("CAM_PASS", "env-pass")
	t.Setenv("CAM_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("NO_PROXY", "true")
	t.Setenv("SCAN_LAST_DAYS", "7")
	t.Setenv("SCAN_FROM_LOCAL_LATEST", "0")
	t.Setenv("HTTP_RETRY_CNT", "6")
	t.Setenv("HTTP_TIMEOUT", "25")
	t.Setenv("DOWNLOAD_DIR", downloadDir)

	cmd, overrides := newConfigOverrideTestCommand()
	if err := cmd.ParseFlags(nil); err != nil {
		t.Fatalf("ParseFlags error: %v", err)
	}

	if err := applyConfigFlagOverrides(cmd.Flags(), overrides); err != nil {
		t.Fatalf("applyConfigFlagOverrides error: %v", err)
	}

	conf, err := config.Init(slog.LevelInfo, false)
	if err != nil {
		t.Fatalf("Init error: %v", err)
	}

	if conf.Host != "env-camera" {
		t.Fatalf("expected env CAM_HOST, got %q", conf.Host)
	}
	if !conf.InsecureTLS {
		t.Fatalf("expected env CAM_INSECURE_SKIP_VERIFY=true")
	}
	if !conf.NoProxy {
		t.Fatalf("expected env NO_PROXY=true")
	}
	if conf.ScanLastDays != 7 {
		t.Fatalf("expected env SCAN_LAST_DAYS, got %d", conf.ScanLastDays)
	}
	if conf.RetryCnt != 6 {
		t.Fatalf("expected env HTTP_RETRY_CNT, got %d", conf.RetryCnt)
	}
	if conf.HttpTimeout != 25 {
		t.Fatalf("expected env HTTP_TIMEOUT, got %d", conf.HttpTimeout)
	}
}

func TestApplyConfigFlagOverridesUsesShortCLIValues(t *testing.T) {
	downloadDir := setBaseConfigEnv(t)
	t.Setenv("CAM_HOST", "env-camera")
	t.Setenv("CAM_USER", "env-user")
	t.Setenv("CAM_PASS", "env-pass")
	t.Setenv("CAM_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("NO_PROXY", "true")
	t.Setenv("SCAN_LAST_DAYS", "9")
	t.Setenv("SCAN_FROM_LOCAL_LATEST", "4")
	t.Setenv("HTTP_RETRY_CNT", "5")
	t.Setenv("HTTP_TIMEOUT", "30")

	cmd, overrides := newConfigOverrideTestCommand()
	if err := cmd.ParseFlags([]string{
		"-d", downloadDir,
		"-H", "short-camera",
		"-u", "short-user",
		"-p", "short-pass",
		"-k=false",
		"-s=0",
		"-l=0",
		"-r=2",
		"-t=15",
		"-P=false",
	}); err != nil {
		t.Fatalf("ParseFlags error: %v", err)
	}

	if err := applyConfigFlagOverrides(cmd.Flags(), overrides); err != nil {
		t.Fatalf("applyConfigFlagOverrides error: %v", err)
	}

	conf, err := config.Init(slog.LevelInfo, false)
	if err != nil {
		t.Fatalf("Init error: %v", err)
	}

	if conf.Host != "short-camera" {
		t.Fatalf("expected short CLI CAM_HOST override, got %q", conf.Host)
	}
	if conf.User != "short-user" {
		t.Fatalf("expected short CLI CAM_USER override, got %q", conf.User)
	}
	if conf.Pass != "short-pass" {
		t.Fatalf("expected short CLI CAM_PASS override, got %q", conf.Pass)
	}
	if conf.InsecureTLS {
		t.Fatalf("expected short CLI CAM_INSECURE_SKIP_VERIFY=false override")
	}
	if conf.NoProxy {
		t.Fatalf("expected short CLI NO_PROXY=false override")
	}
	if conf.ScanLastDays != 0 {
		t.Fatalf("expected short CLI SCAN_LAST_DAYS override, got %d", conf.ScanLastDays)
	}
	if conf.ScanLocal != 0 {
		t.Fatalf("expected short CLI SCAN_FROM_LOCAL_LATEST override, got %d", conf.ScanLocal)
	}
	if conf.RetryCnt != 2 {
		t.Fatalf("expected short CLI HTTP_RETRY_CNT override, got %d", conf.RetryCnt)
	}
	if conf.HttpTimeout != 15 {
		t.Fatalf("expected short CLI HTTP_TIMEOUT override, got %d", conf.HttpTimeout)
	}
}

func TestConfigFlagsAppearInHelp(t *testing.T) {
	cmd := newRootCmd()
	help := cmd.UsageString()
	for _, expected := range []string{
		"-d, --download-dir string",
		"-H, --cam-host string",
		"-u, --cam-user string",
		"-p, --cam-pass string",
		"-k, --cam-insecure-skip-verify",
		"-s, --scan-last-days int",
		"-l, --scan-from-local-latest int",
		"-r, --http-retry-cnt int",
		"-t, --http-timeout int",
		"-P, --no-proxy",
	} {
		if !strings.Contains(help, expected) {
			t.Fatalf("expected help to contain %q, got:\n%s", expected, help)
		}
	}
}

func newConfigOverrideTestCommand() (*cobra.Command, *configFlagOverrides) {
	cmd := &cobra.Command{Use: "test"}
	overrides := &configFlagOverrides{}
	addConfigFlags(cmd, overrides)
	return cmd, overrides
}

func setBaseConfigEnv(t *testing.T) string {
	t.Helper()

	downloadDir := t.TempDir()
	t.Setenv("DOWNLOAD_DIR", downloadDir)
	t.Setenv("CAM_HOST", "base-camera")
	t.Setenv("CAM_USER", "base-user")
	t.Setenv("CAM_PASS", "base-pass")
	t.Setenv("CAM_INSECURE_SKIP_VERIFY", "false")
	t.Setenv("NO_PROXY", "false")
	t.Setenv("SCAN_LAST_DAYS", "0")
	t.Setenv("SCAN_FROM_LOCAL_LATEST", "0")
	t.Setenv("HTTP_RETRY_CNT", "3")
	t.Setenv("HTTP_TIMEOUT", "120")

	return downloadDir
}
