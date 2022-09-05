package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var repoPath string
var rootCmd = &cobra.Command{
	Use:   "owl",
	Short: "Owl Blogs is a not so static blog generator",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&repoPath, "repo", ".", "Path to the repository to use.")
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "Username")

}

func main() {
	Execute()
}
