package main

import (
	"h4kor/kiss-social"
	"os"
)

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
		kiss.CreateNewUser("users", args[1])
	default:
		println("Invalid command")
	}
}
