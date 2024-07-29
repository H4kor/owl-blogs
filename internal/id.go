package internal

import (
	"fmt"
	"regexp"
	"strings"
)

func TurnIntoId(name string, availCheck func(string) bool) string {
	// try to find a good ID
	m := regexp.MustCompile(`[^a-z0-9-]`)
	prefix := m.ReplaceAllString(strings.ToLower(name), "-")
	id := prefix + ""
	counter := 0
	for {
		avail := availCheck(id)
		if !avail {
			counter += 1
			id = fmt.Sprintf("%s-%d", prefix, counter)
		} else {
			break
		}
	}
	return id
}
