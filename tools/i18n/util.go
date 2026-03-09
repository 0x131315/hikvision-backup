package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func hashString(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum[:])
}

func isQuotaExceeded(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "quota exceeded")
}

func htmlUnescape(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
	)
	return replacer.Replace(s)
}
