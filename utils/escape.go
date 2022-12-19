package utils

import "strings"

func FixEscape(str string) string {
	fixed := strings.ReplaceAll(str, "'", "''")
	return fixed
}