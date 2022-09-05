package main

import (
	web "h4kor/owl-blogs/cmd/owl/web"

	"github.com/spf13/cobra"
)

var port int
var unsafe bool
var user string

func init() {
	rootCmd.AddCommand(webCmd)

	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "Port to use")
	rootCmd.PersistentFlags().BoolVar(&unsafe, "unsafe", false, "Allow unsafe html")
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "Start server in single user mode.")
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Long:  `Start the web server`,
	Run: func(cmd *cobra.Command, args []string) {
		web.StartServer(repoPath, port, unsafe, user)
	},
}
