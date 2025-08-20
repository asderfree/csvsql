package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// IsChineseChar checks if a character is Chinese
func IsChineseChar(r rune) bool {
	return unicode.Is(unicode.Han, r)
}

// ContainsChinese checks if a string contains Chinese characters
func ContainsChinese(s string) bool {
	for _, r := range s {
		if IsChineseChar(r) {
			return true
		}
	}
	return false
}

// SanitizeTableName sanitizes a table name to be valid SQL
func SanitizeTableName(name string) string {
	// Remove leading special characters
	reg := regexp.MustCompile("^[^a-zA-Z0-9_]+")
	sanitized := reg.ReplaceAllString(name, "")

	// Replace invalid characters with underscores
	reg = regexp.MustCompile("[^a-zA-Z0-9_]+")
	sanitized = reg.ReplaceAllString(sanitized, "")

	// Prepend underscore if table name starts with a number
	if matched, _ := regexp.MatchString("^[0-9]", sanitized); matched {
		sanitized = "_" + sanitized
	}

	return sanitized
}

// SanitizeColumnName sanitizes a column name for SQL use
func SanitizeColumnName(name string) string {
	// Replace spaces with underscores
	sanitized := strings.ReplaceAll(name, " ", "_")

	// Replace invalid characters
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	return reg.ReplaceAllString(sanitized, "")
}
