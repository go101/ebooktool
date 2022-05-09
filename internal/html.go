package internal

import (
	"regexp"
)

var regexpInvalidIdentifier = regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
var regexpInvalidIdentifierExcludeDots = regexp.MustCompile(`[^a-zA-Z0-9_\-\.]`)

func ValidateIdentifier(id string, keepDots bool) string {
	if keepDots {
		return regexpInvalidIdentifierExcludeDots.ReplaceAllString(id, "-")
	}
	return regexpInvalidIdentifier.ReplaceAllString(id, "-")
}
