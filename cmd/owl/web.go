package main

import (
	owlblogs "owl-blogs"
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

var bindAddr string

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVarP(&bindAddr, "bind", "b", "localhost:3000", "Address to bind to")

}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		owlblogs.App(db).Run(bindAddr)
	},
}
