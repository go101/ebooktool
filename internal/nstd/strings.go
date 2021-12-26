package nstd

import (
	"strings"
)

type String string

func (s String) String() string {
	return string(s)
}

func (s String) Index(sub string) int {
	return strings.Index(string(s), sub)
}

func (s String) ToLower() String {
	return String(strings.ToLower(string(s)))
}

func (s String) HasPrefix(prefix string) bool {
	return strings.HasPrefix(string(s), prefix)
}

func (s String) HasSuffix(suffix string) bool {
	return strings.HasSuffix(string(s), suffix)
}

func (s String) TrimSuffix(suffix string) String {
	return String(strings.TrimSuffix(string(s), suffix))
}

func (s String) ReplaceSuffix(suffix, with string) String {
	if s.HasSuffix(suffix) {
		return s[:len(s)-len(suffix)] + String(with)
	}
	return s
}

// Note, the SplitN method returns []string instead of []String.
func (s String) SplitN(sep string, n int) []string {
	return strings.SplitN(string(s), sep, n)
}

func (s String) TrimSpace() String {
	return String(strings.TrimSpace(string(s)))
}

func (s String) ReplaceAll(from, to string) String {
	return String(strings.ReplaceAll(string(s), from, to))
}

func (s String) LastIndex(sub string) int {
	return strings.LastIndex(string(s), sub)
}
