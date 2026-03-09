package main

import (
	"fmt"
	"os"
)

func run(checkOnly, initMode, force bool) error {
	source, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", sourceFile, err)
	}

	blocks, seps := splitBlocks(string(source))
	sourceHash := hashString(string(source))

	for _, lang := range langs {
		if err := processLanguage(lang, blocks, seps, sourceHash, checkOnly, initMode, force); err != nil {
			return err
		}
	}

	return nil
}
