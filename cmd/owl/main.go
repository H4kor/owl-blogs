package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var DbPath string

var rootCmd = &cobra.Command{
	Use:   "owl",
	Short: "Owl Blogs is a not so static blog generator",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&DbPath, "file", "f", "owlblogs.db", "Path to blog file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
