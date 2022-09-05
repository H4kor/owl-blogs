package main

import (
	web "h4kor/owl-blogs/cmd/owl/web"

	"github.com/spf13/cobra"
)

var port int

func init() {
	rootCmd.AddCommand(webCmd)

	webCmd.PersistentFlags().IntVar(&port, "port", 8080, "Port to use")
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		web.StartServer(repoPath, port)
	},
}
