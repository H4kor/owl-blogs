package main

import (
	"h4kor/owl-blogs"

	"github.com/spf13/cobra"
)

var domain string
var singleUser string
var unsafe bool

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringVar(&domain, "domain", "http://localhost:8080", "Domain to use")
	initCmd.PersistentFlags().StringVar(&singleUser, "single-user", "", "Use single user mode with given username")
	initCmd.PersistentFlags().BoolVar(&unsafe, "unsafe", false, "Allow raw html")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a new repository",
	Long:  `Creates a new repository`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := owl.CreateRepository(repoPath, owl.RepoConfig{
			Domain:       domain,
			SingleUser:   singleUser,
			AllowRawHtml: unsafe,
		})
		if err != nil {
			println("Error creating repository: ", err.Error())
		} else {
			println("Repository created: ", repoPath)
		}

	},
}
