package main

import (
	"h4kor/owl-blogs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webmentionCmd)
}

var webmentionCmd = &cobra.Command{
	Use:   "webmention",
	Short: "Send webmentions for posts, optionally for a specific user",
	Long:  `Send webmentions for posts, optionally for a specific user`,
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := owl.OpenRepository(repoPath)
		if err != nil {
			println("Error opening repository: ", err.Error())
			return
		}

		var users []owl.User
		if user == "" {
			// send webmentions for all users
			users, err = repo.Users()
			if err != nil {
				println("Error getting users: ", err.Error())
				return
			}
		} else {
			// send webmentions for a specific user
			user, err := repo.GetUser(user)
			users = append(users, user)
			if err != nil {
				println("Error getting user: ", err.Error())
				return
			}
		}

		for _, user := range users {
			posts, err := user.Posts()
			if err != nil {
				println("Error getting posts: ", err.Error())
			}

			for _, post := range posts {
				println("Webmentions for post: ", post.Title())

				err := post.ScanForLinks()
				if err != nil {
					println("Error scanning post for links: ", err.Error())
					continue
				}

				webmentions := post.OutgoingWebmentions()
				println("Found ", len(webmentions), " links")
				for _, webmention := range webmentions {
					err = post.SendWebmention(webmention)
					if err != nil {
						println("Error sending webmentions: ", err.Error())
					} else {
						println("Webmention sent to ", webmention.Target)
					}
				}
			}
		}
	},
}
