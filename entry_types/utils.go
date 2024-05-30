package entrytypes

import (
	"regexp"
	"strings"
)

func extractTags(content string) []string {
	r := regexp.MustCompile("#[a-zA-Z0-9_]+")
	matches := r.FindAllString(string(content), -1)
	tags := make([]string, 0)
	for _, hashtag := range matches {
		tag, _ := strings.CutPrefix(hashtag, "#")
		tags = append(tags, tag)
	}
	return tags

}
