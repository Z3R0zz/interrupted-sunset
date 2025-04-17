package utils

import (
	"regexp"
	"strings"
)

var invalidChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

func SanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = invalidChars.ReplaceAllString(name, "_")

	if len(name) > 128 {
		name = name[:128]
	}

	return name
}
