package entrytypes

import (
	"strings"
	"unicode"
)

func endOfHashtag(r rune) bool {
	return !unicode.IsLetter(r) &&
		!unicode.IsDigit(r) &&
		r != '_' && r != '-' && r != '/'
}

func extractTags(content string) []string {
	tagMap := make(map[string]bool)
	start := -1
	for i, c := range content {
		if start != -1 {
			if endOfHashtag(c) {
				if i != start {
					tag := strings.ToLower(content[start:i])
					tagMap[tag] = true
				}
				start = -1
			}
		} else {
			if c == rune('#') && (i == 0 || content[i-1] == ' ' || content[i-1] == '\n') {
				start = i + 1
				continue
			}
		}
	}
	if start != -1 && len(content)+1 > start {
		tag := strings.ToLower(content[start:])
		tagMap[tag] = true
	}
	tags := make([]string, 0, len(tagMap))
	for t := range tagMap {
		tags = append(tags, t)
	}
	return tags
}
