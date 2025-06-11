package main

import (
	"fmt"
	"github.com/0x131315/hikvision-backup/internal/app"
	"os"
)

// Эти переменные будут "заполнены" во время сборки
var (
	version   = "development" // Версия сборки
	commit    = "none"        // Хэш коммита
	buildDate = "unknown"     // Дата сборки
)

func main() {
	handlePrintVersionCmd()
	app.DownloadVideos()
}

func handlePrintVersionCmd() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("Version:   %s\n", version)
		fmt.Printf("Commit:    %s\n", commit)
		fmt.Printf("BuildDate: %s\n", buildDate)
		os.Exit(0)
	}
}
