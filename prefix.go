package log

import (
	"strings"
)

const (
	suffix = ": "
)

func Prefixed(prefix, name string) string {
	if prefix != "" {
		return prefix + "." + name
	}
	return name
}

func PrefixedWithSuffix(prefix, name string) string {
	return Prefixed(strings.TrimSuffix(prefix, suffix), name) + suffix
}
