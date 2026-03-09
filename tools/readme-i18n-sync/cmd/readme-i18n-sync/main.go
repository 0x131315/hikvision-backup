package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0x131315/readme-i18n-sync/internal/syncer"
)

func main() {
	checkOnly := flag.Bool("check", false, "check translations without updating")
	initMode := flag.Bool("init", false, "bootstrap translations using source text when missing")
	force := flag.Bool("force", false, "retranslate all blocks and overwrite existing translations")
	source := flag.String("source", "README.md", "source README path")
	i18nDir := flag.String("i18n-dir", "i18n", "translated README output directory")
	tmDir := flag.String("tm-dir", "", "translation memory directory (default: <i18n-dir>/tm)")
	flag.Parse()

	syncer.SetPaths(*source, *i18nDir, *tmDir)

	if err := syncer.Run(*checkOnly, *initMode, *force); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
