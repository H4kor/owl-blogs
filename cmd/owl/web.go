package main

import (
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

var port string
var host string

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVarP(&port, "port", "p", "3000", "Port to listen on")
	webCmd.Flags().StringVarP(&host, "host", "b", "localhost", "Address to listen on")

}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		App(db).Run(host, port)
	},
}
