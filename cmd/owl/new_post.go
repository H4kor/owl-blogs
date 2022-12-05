package main

import (
	"h4kor/owl-blogs"

	"github.com/spf13/cobra"
)

var postTitle string

func init() {
	rootCmd.AddCommand(newPostCmd)
	newPostCmd.PersistentFlags().StringVar(&postTitle, "title", "", "Post title")
}

var newPostCmd = &cobra.Command{
	Use:   "new-post",
	Short: "Creates a new post",
	Long:  `Creates a new post`,
	Run: func(cmd *cobra.Command, args []string) {
		if user == "" {
			println("Username is required")
			return
		}

		if postTitle == "" {
			println("Post title is required")
			return
		}

		repo, err := owl.OpenRepository(repoPath)
		if err != nil {
			println("Error opening repository: ", err.Error())
			return
		}

		user, err := repo.GetUser(user)
		if err != nil {
			println("Error getting user: ", err.Error())
			return
		}

		post, err := user.CreateNewPost(owl.PostMeta{Type: "article", Title: postTitle, Draft: true}, "")
		if err != nil {
			println("Error creating post: ", err.Error())
		} else {
			println("Post created: ", postTitle)
			println("Edit: ", post.ContentFile())
		}
	},
}
