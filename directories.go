package kiss

import "os"

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
