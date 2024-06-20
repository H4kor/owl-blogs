package main

import (
	owlblogs "owl-blogs"
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetPasswordCmd)

	resetPasswordCmd.Flags().StringVarP(&user, "user", "u", "", "The user name")
	resetPasswordCmd.MarkFlagRequired("user")
	resetPasswordCmd.Flags().StringVarP(&password, "password", "p", "", "The new password")
	resetPasswordCmd.MarkFlagRequired("password")
}

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Resets the password of an author",
	Long:  `Resets the password of an author`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		owlblogs.App(db).AuthorService.Create(user, password)
	},
}
