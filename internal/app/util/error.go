package util

import (
	"fmt"
	"os"
	"runtime/debug"
)

func FatalError(msg string, err ...error) {
	fmt.Printf("%s: %s\n trace: %s \n", msg, err, debug.Stack())
	os.Exit(1)
}
