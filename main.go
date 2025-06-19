package main

import (
	"github.com/0x131315/hikvision-backup/internal/app"
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		printVersion()
		os.Exit(0)
	}

	level := slog.LevelInfo
	if len(os.Args) > 1 && os.Args[1] == "-vv" {
		level = slog.LevelDebug
	}
	initLogger(level)

	app.DownloadVideos()
}
