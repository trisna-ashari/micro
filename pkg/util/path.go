package util

import "strings"

// MakePathWithPrefix is a function uses to add prefix into given path.
func MakePathWithPrefix(prefix string, path string) string {
	var str []string
	if len(prefix) > 0 {
		str = append(str, prefix)
	}

	str = append(str, path)

	return strings.Join(str, "/")
}
