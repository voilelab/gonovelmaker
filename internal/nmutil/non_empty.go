package nmutil

import (
	"strings"
)

func FirstNonEmptyString(strs ...string) string {
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str != "" {
			return str
		}
	}
	return ""
}

func FirstNonZero[T comparable](vals ...T) T {
	var zero T
	for _, val := range vals {
		if val != zero {
			return val
		}
	}
	return zero
}
