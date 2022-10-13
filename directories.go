package owl

import (
	"os"
	"strings"
)

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// lists all files/dirs in a directory, not recursive
func listDir(path string) []string {
	dir, _ := os.Open(path)
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)
	return files
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func toDirectoryName(name string) string {
	name = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	// remove all non-alphanumeric characters
	name = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r
		}
		if r >= '0' && r <= '9' {
			return r
		}
		if r == '-' {
			return r
		}
		return -1
	}, name)
	return name
}
