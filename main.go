package main

import (
	"os"
	"path"
)

func CreateNewUser(repo string, name string) {
	// creates repo folder if it doesn't exist
	os.Mkdir(repo, 0755)

	// creates repo/name folder if it doesn't exist
	user_dir := path.Join(repo, name)
	os.Mkdir(user_dir, 0755)
	os.Mkdir(path.Join(user_dir, "meta"), 0755)
	os.WriteFile(path.Join(user_dir, "meta", "VERSION"), []byte("0.0.1"), 0644)
	os.WriteFile(path.Join(user_dir, "meta", "base.html"), []byte("<html><body><{{content}}/body></html>"), 0644)
}

func main() {
	println("KISS Social")
	println("Commands")
	println("new <name> - Creates a new user")

	args := os.Args[1:]
	if len(args) == 0 {
		println("No command given")
		return
	}

	switch args[0] {
	case "new":
		if len(args) != 2 {
			println("Invalid number of arguments")
			return
		}
		CreateNewUser("users", args[1])
	default:
		println("Invalid command")
	}
}