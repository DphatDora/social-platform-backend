package util

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// NormalizeString removes accents and diacritics from a string,
// converts it to lowercase, and trims leading/trailing spaces.
func NormalizeString(s string) string {
	// Decompose Unicode characters (e.g., "ấ" → "a" + "̂" + "́")
	t := norm.NFD.String(s)

	sb := strings.Builder{}
	for _, r := range t {
		// Skip non-spacing marks (diacritics)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		// Keep only base characters
		sb.WriteRune(r)
	}

	// Convert to lowercase and trim spaces
	return strings.ToLower(strings.TrimSpace(sb.String()))
}

// BuildSearchPattern splits keywords and creates a search pattern for SQL LIKE
func BuildSearchPattern(keyword string) []string {
	normalized := NormalizeString(keyword)
	parts := strings.Fields(normalized) // split by space
	patterns := make([]string, len(parts))

	for i, p := range parts {
		patterns[i] = "%" + p + "%"
	}
	return patterns
}
