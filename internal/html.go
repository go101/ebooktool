package internal

import (
	"regexp"
)

var regexpInvalidIdentifier = regexp.MustCompile(`[^a-zA-Z0-9_\-]`)

func ValidateIdentifier(id string) string {
	return regexpInvalidIdentifier.ReplaceAllString(id, "-")
}
