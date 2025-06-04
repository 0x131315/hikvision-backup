package util

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

func BuildProgressBar(size int, sizeType string) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		size,
		progressbar.OptionSetItsString(sizeType),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowTotalBytes(true),
		progressbar.OptionUseIECUnits(true),
		progressbar.OptionSetWidth(80),
		progressbar.OptionThrottle(100*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(9),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
	)
}
