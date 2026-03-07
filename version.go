package main

import (
	"fmt"
)

// These variables are populated at build time.
var (
	version   = "development" // Build version
	commit    = "none"        // Commit hash
	buildDate = "unknown"     // Build date
)

func printVersion() {
	fmt.Printf("Version:   %s\n", version)
	fmt.Printf("Commit:    %s\n", commit)
	fmt.Printf("BuildDate: %s\n", buildDate)
	fmt.Printf("Source: %s\n", "https://github.com/0x131315/hikvision-backup")
}
