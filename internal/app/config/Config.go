package config

import (
	"github.com/0x131315/hikvision-backup/internal/app/util"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

type Config struct {
	Host         string
	User         string
	Pass         string
	NoProxy      bool
	DownloadDir  string
	ScanLastDays int
}

var config *Config

func init() {
	conf := buildConfig()
	config = &conf
}

func Get() Config {
	return *config
}

func buildConfig() Config {
	host := os.Getenv("CAM_HOST")
	if host == "" {
		util.FatalError("CAM_HOST environment variable not set")
	}

	user := os.Getenv("CAM_USER")
	if user == "" {
		util.FatalError("CAM_USER environment variable not set")
	}

	pass := os.Getenv("CAM_PASS")

	noProxy := os.Getenv("NO_PROXY")
	if noProxy == "" {
		noProxy = "false"
	}

	downloadDir := os.Getenv("DOWNLOAD_DIR")
	if downloadDir == "" {
		util.FatalError("DOWNLOAD_DIR environment variable not set")
	}

	scanLastDays := os.Getenv("SCAN_LAST_DAYS")
	if scanLastDays == "" {
		util.FatalError("SCAN_LAST_DAYS environment variable not set")
	}

	lastDays, err := strconv.Atoi(scanLastDays)
	if err != nil {
		util.FatalError("SCAN_LAST_DAYS environment variable not numeric")
	}
	if lastDays < 0 {
		lastDays = lastDays * -1
	}
	if lastDays < 1 {
		util.FatalError("SCAN_LAST_DAYS environment variable smaller than 1")
	}

	conf := Config{
		Host:         host,
		User:         user,
		Pass:         pass,
		NoProxy:      noProxy == "true",
		DownloadDir:  downloadDir,
		ScanLastDays: lastDays,
	}

	return conf
}
