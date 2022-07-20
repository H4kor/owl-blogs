package kiss

import "os"

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func listDir(path string) []string {
	dir, _ := os.Open(path)
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)
	return files
}
