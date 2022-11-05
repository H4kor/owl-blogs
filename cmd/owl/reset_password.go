package main

import (
	"fmt"
	"h4kor/owl-blogs"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetPasswordCmd)
}

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset the password for a user",
	Long:  `Reset the password for a user`,
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

		user, err := repo.GetUser(user)
		if err != nil {
			println("Error getting user: ", err.Error())
			return
		}

		// generate a random password and print it
		password := owl.GenerateRandomString(16)
		user.ResetPassword(password)

		fmt.Println("User:         ", user.Name())
		fmt.Println("New Password: ", password)

	},
}
