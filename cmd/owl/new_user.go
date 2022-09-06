package main

import (
	"h4kor/owl-blogs"

	"github.com/spf13/cobra"
)

var user string

func init() {
	rootCmd.AddCommand(newUserCmd)
}

var newUserCmd = &cobra.Command{
	Use:   "new-user",
	Short: "Creates a new user",
	Long:  `Creates a new user`,
	Run: func(cmd *cobra.Command, args []string) {
		if user == "" {
			println("Username is required")
			return
		}

		repo, err := owl.OpenRepository(repoPath)
		if err != nil {
			println("Error opening repository: ", err.Error())
			return
		}

		_, err = repo.CreateUser(user)
		if err != nil {
			println("Error creating user: ", err.Error())
		} else {
			println("User created: ", user)
		}
	},
}
