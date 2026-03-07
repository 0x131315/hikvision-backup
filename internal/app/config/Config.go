package config

import (
	"log/slog"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Host         string
	User         string
	Pass         string
	NoProxy      bool
	DownloadDir  string
	ScanLastDays int
	ScanLocal    int
	RetryCnt     int
	HttpTimeout  int
	LogLvl       slog.Level
	LogHttp      bool
}

var config *Config

func Init(logLvl slog.Level, logHttp bool) Config {
	conf := buildConfig(logLvl, logHttp)
	config = &conf

	return conf
}

func Get() Config {
	return *config
}

func buildConfig(logLvl slog.Level, logHttp bool) Config {
	host := os.Getenv("CAM_HOST")
	if host == "" {
		slog.Error("CAM_HOST environment variable not set")
		os.Exit(1)
	}

	user := os.Getenv("CAM_USER")
	if user == "" {
		slog.Error("CAM_USER environment variable not set")
		os.Exit(1)
	}

	pass := os.Getenv("CAM_PASS")

	noProxy := os.Getenv("NO_PROXY")
	if noProxy == "" {
		noProxy = "false"
	}

	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		slog.Error("DOWNLOAD_DIR environment variable not set")
		os.Exit(1)
	}

	scanLastDays := os.Getenv("SCAN_LAST_DAYS")
	if scanLastDays == "" {
		slog.Debug("SCAN_LAST_DAYS environment variable not set")
		scanLastDays = "0"
	}

	lastDays, err := strconv.Atoi(scanLastDays)
	if err != nil {
		slog.Error("SCAN_LAST_DAYS environment variable not numeric")
		os.Exit(1)
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
		slog.Error("SCAN_FROM_LOCAL_LATEST environment variable not numeric")
		os.Exit(1)
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
		slog.Error("HTTP_RETRY_CNT environment variable not numeric")
		os.Exit(1)
	}
	if retryCnt < 0 {
		slog.Error("HTTP_RETRY_CNT environment variable smaller than 0")
		os.Exit(1)
	}

	envHttpTimeout := os.Getenv("HTTP_TIMEOUT")
	if envHttpTimeout == "" {
		slog.Debug("HTTP_TIMEOUT environment variable not set")
		envHttpTimeout = "120"
	}
	httpTimeout, err := strconv.Atoi(envHttpTimeout)
	if err != nil {
		slog.Error("HTTP_TIMEOUT environment variable not numeric")
		os.Exit(1)
	}
	if httpTimeout < 0 {
		slog.Error("HTTP_TIMEOUT environment variable smaller than 0")
		os.Exit(1)
	}

	conf := Config{
		Host:         host,
		User:         user,
		Pass:         pass,
		NoProxy:      noProxy == "true",
		DownloadDir:  downloadDir,
		ScanLastDays: lastDays,
		ScanLocal:    scanLocalDays,
		RetryCnt:     retryCnt,
		HttpTimeout:  httpTimeout,
		LogLvl:       logLvl,
		LogHttp:      logHttp,
	}

	return conf
}
