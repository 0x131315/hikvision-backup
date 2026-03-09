package config

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host         string
	BaseURL      string
	User         string
	Pass         string
	NoProxy      bool
	InsecureTLS  bool
	DownloadDir  string
	ScanLastDays int
	ScanLocal    int
	RetryCnt     int
	HttpTimeout  int
	LogLvl       slog.Level
	LogHttp      bool
}

func Init(logLvl slog.Level, logHttp bool) (Config, error) {
	return buildConfig(logLvl, logHttp)
}

func buildConfig(logLvl slog.Level, logHttp bool) (Config, error) {
	host := os.Getenv("CAM_HOST")
	if host == "" {
		return Config{}, fmt.Errorf("CAM_HOST environment variable not set")
	}
	baseURL, err := buildBaseURL(host)
	if err != nil {
		return Config{}, err
	}

	user := os.Getenv("CAM_USER")
	if user == "" {
		return Config{}, fmt.Errorf("CAM_USER environment variable not set")
	}

	pass := os.Getenv("CAM_PASS")

	noProxy := os.Getenv("NO_PROXY")
	if noProxy == "" {
		noProxy = "false"
	}
	insecureTLS := os.Getenv("CAM_INSECURE_SKIP_VERIFY")
	if insecureTLS == "" {
		insecureTLS = "false"
	}
	insecureSkipVerify, err := strconv.ParseBool(insecureTLS)
	if err != nil {
		return Config{}, fmt.Errorf("CAM_INSECURE_SKIP_VERIFY environment variable must be true/false")
	}

	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		return Config{}, fmt.Errorf("DOWNLOAD_DIR environment variable not set")
	}

	scanLastDays := os.Getenv("SCAN_LAST_DAYS")
	if scanLastDays == "" {
		slog.Debug("SCAN_LAST_DAYS environment variable not set")
		scanLastDays = "0"
	}

	lastDays, err := strconv.Atoi(scanLastDays)
	if err != nil {
		return Config{}, fmt.Errorf("SCAN_LAST_DAYS environment variable not numeric")
	}
	if lastDays < 0 {
		lastDays = lastDays * -1
	}

	scanLocalStr := os.Getenv("SCAN_FROM_LOCAL_LATEST")
	if scanLocalStr == "" {
		scanLocalStr = "0"
	}
	scanLocalDays, err := strconv.Atoi(scanLocalStr)
	if err != nil {
		return Config{}, fmt.Errorf("SCAN_FROM_LOCAL_LATEST environment variable not numeric")
	}
	if scanLocalDays < 0 {
		scanLocalDays = scanLocalDays * -1
	}
	if scanLocalDays > 0 {
		lastDays = adjustScanLastDaysFromLocal(downloadDir, lastDays, scanLocalDays)
	}

	envRetryCnt := os.Getenv("HTTP_RETRY_CNT")
	if envRetryCnt == "" {
		slog.Debug("HTTP_RETRY_CNT environment variable not set")
		envRetryCnt = "3"
	}
	retryCnt, err := strconv.Atoi(envRetryCnt)
	if err != nil {
		return Config{}, fmt.Errorf("HTTP_RETRY_CNT environment variable not numeric")
	}
	if retryCnt < 0 {
		return Config{}, fmt.Errorf("HTTP_RETRY_CNT environment variable smaller than 0")
	}

	envHttpTimeout := os.Getenv("HTTP_TIMEOUT")
	if envHttpTimeout == "" {
		slog.Debug("HTTP_TIMEOUT environment variable not set")
		envHttpTimeout = "120"
	}
	httpTimeout, err := strconv.Atoi(envHttpTimeout)
	if err != nil {
		return Config{}, fmt.Errorf("HTTP_TIMEOUT environment variable not numeric")
	}
	if httpTimeout < 0 {
		return Config{}, fmt.Errorf("HTTP_TIMEOUT environment variable smaller than 0")
	}

	conf := Config{
		Host:         host,
		BaseURL:      baseURL,
		User:         user,
		Pass:         pass,
		NoProxy:      noProxy == "true",
		InsecureTLS:  insecureSkipVerify,
		DownloadDir:  downloadDir,
		ScanLastDays: lastDays,
		ScanLocal:    scanLocalDays,
		RetryCnt:     retryCnt,
		HttpTimeout:  httpTimeout,
		LogLvl:       logLvl,
		LogHttp:      logHttp,
	}

	return conf, nil
}

func buildBaseURL(host string) (string, error) {
	host = strings.TrimSpace(host)
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		u, err := url.Parse(host)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return "", fmt.Errorf("CAM_HOST must be a valid host or URL")
		}
		return strings.TrimRight(host, "/"), nil
	}

	return "https://" + host, nil
}
