package main

import (
	"h4kor/owl-blogs"
	"os"
)

func main() {
	println("owl blogs")
	println("Commands")
	println("init <repo> - Creates a new repository")
	println("<repo> new-user <name> - Creates a new user")
	println("<repo> new-post <user> <title> - Creates a new post")

	if len(os.Args) < 3 {
		println("Please specify a repository and command")
		os.Exit(1)
	}

	if os.Args[1] == "init" {
		repoName := os.Args[2]
		_, err := owl.CreateRepository(repoName)
		if err != nil {
			println("Error creating repository: ", err.Error())
		}
		println("Repository created: ", repoName)
		os.Exit(0)
	}

	repoName := os.Args[1]
	repo, err := owl.OpenRepository(repoName)
	if err != nil {
		println("Error opening repository: ", err.Error())
		os.Exit(1)
	}
	switch os.Args[2] {
	case "new-user":
		if len(os.Args) < 4 {
			println("Please specify a user name")
			os.Exit(1)
		}
		userName := os.Args[3]
		user, err := repo.CreateUser(userName)
		if err != nil {
			println("Error creating user: ", err.Error())
			os.Exit(1)
		}
		println("User created: ", user.Name())
	case "new-post":
		if len(os.Args) < 5 {
			println("Please specify a user name and a title")
			os.Exit(1)
		}
		userName := os.Args[3]
		user, err := repo.GetUser(userName)
		if err != nil {
			println("Error finding user: ", err.Error())
			os.Exit(1)
		}
		title := os.Args[4]
		post, err := user.CreateNewPost(title)
		if err != nil {
			println("Error creating post: ", err.Error())
			os.Exit(1)
		}
		println("Post created: ", post.Title())
	default:
		println("Unknown command: ", os.Args[2])
		os.Exit(1)
	}
}
