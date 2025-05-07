package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// RemoveAccents removes Vietnamese accents and returns ASCII string.
func RemoveAccents(input string) string {
	t := norm.NFD.String(input)
	var output strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		output.WriteRune(r)
	}
	return output.String()
}

// Slugify generates a slug from product name.
func Slugify(name string) string {
	// Remove accents
	name = RemoveAccents(name)

	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace non-alphanumeric characters with hyphens
	reg, _ := regexp.Compile("[^a-z0-9]+")
	name = reg.ReplaceAllString(name, "-")

	// Trim hyphens
	name = strings.Trim(name, "-")

	return name
}
