package owl

import (
	"os"
	"path/filepath"
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

// recursive list of all files in a directory
func walkDir(path string) []string {
	files := make([]string, 0)
	filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, subPath[len(path)+1:])
		return nil
	})
	return files
}
