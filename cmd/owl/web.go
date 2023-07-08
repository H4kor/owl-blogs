package main

import (
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		App(db).Run()
	},
}
