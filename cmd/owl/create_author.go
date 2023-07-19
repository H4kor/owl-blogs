package main

import (
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

var user string
var password string

func init() {
	rootCmd.AddCommand(newAuthorCmd)

	newAuthorCmd.Flags().StringVarP(&user, "user", "u", "", "The user name")
	newAuthorCmd.MarkFlagRequired("user")
	newAuthorCmd.Flags().StringVarP(&password, "password", "p", "", "The password")
	newAuthorCmd.MarkFlagRequired("password")
}

var newAuthorCmd = &cobra.Command{
	Use:   "new-author",
	Short: "Creates a new author",
	Long:  `Creates a new author`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		App(db).AuthorService.Create(user, password)
	},
}
