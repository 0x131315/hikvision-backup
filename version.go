package main

import (
	"fmt"
)

// Эти переменные будут "заполнены" во время сборки
var (
	version   = "development" // Версия сборки
	commit    = "none"        // Хэш коммита
	buildDate = "unknown"     // Дата сборки
)

func printVersion() {
	fmt.Printf("Version:   %s\n", version)
	fmt.Printf("Commit:    %s\n", commit)
	fmt.Printf("BuildDate: %s\n", buildDate)
}
