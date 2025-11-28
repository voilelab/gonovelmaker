package main

import "strings"

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "!", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	// Keep alphanumeric characters, underscore, and Unicode letters (including Chinese, Japanese, etc.)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r > 127 {
			result.WriteRune(r)
		}
	}
	return result.String()
}
