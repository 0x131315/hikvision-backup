package main

import (
	"github.com/0x131315/hikvision-backup/internal/app"
	"os"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		printVersion()
		os.Exit(0)
	}
	app.DownloadVideos()
}
