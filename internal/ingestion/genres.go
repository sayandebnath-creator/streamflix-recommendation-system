package ingestion

import (
	"regexp"
	"strings"
)

var genreRegex = regexp.MustCompile(`'name':\s*'([^']+)'`)

func ParseGenres(raw string) string {
	matches := genreRegex.FindAllStringSubmatch(raw, -1)

	var genres []string

	for _, match := range matches {
		if len(match) > 1 {
			genres = append(genres, match[1])
		}
	}

	return strings.Join(genres, ",")
}