package util

import (
	"fmt"
	"math"
)

func FormatFileSize(s int) string {
	size, suff := humanizeBytes(float64(s), true)
	return fmt.Sprintf("%s%s", size, suff)
}

func humanizeBytes(s float64, iec bool) (string, string) {
	sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
	base := 1000.0

	if iec {
		sizes = []string{" B", " KiB", " MiB", " GiB", " TiB", " PiB", " EiB"}
		base = 1024.0
	}

	if s < 10 {
		return fmt.Sprintf("%2.0f", s), sizes[0]
	}
	e := math.Floor(logn(s, base))
	suffix := sizes[int(e)]
	val := math.Floor(s/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f, val), suffix
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
