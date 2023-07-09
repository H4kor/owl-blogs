package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&user, "path", "p", "", "Path to the user folder")
	importCmd.MarkFlagRequired("path")
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data from v1",
	Long:  `Import data from v1`,
	Run: func(cmd *cobra.Command, args []string) {
		// db := infra.NewSqliteDB(DbPath)
		// App(db).ImportV1()

		// TODO: Implement this
		// For each folder in the user folder
		// Map to entry types
		// Convert and save
		// Import Binary files
	},
}
