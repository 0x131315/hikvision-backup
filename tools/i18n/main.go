package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	checkOnly := flag.Bool("check", false, "check translations without updating")
	initMode := flag.Bool("init", false, "bootstrap translations using source text when missing")
	force := flag.Bool("force", false, "retranslate all blocks and overwrite existing translations")
	flag.Parse()

	if err := run(*checkOnly, *initMode, *force); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
