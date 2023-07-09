package main

import (
	"fmt"
	"owl-blogs/domain/model"
	"owl-blogs/importer"
	"owl-blogs/infra"

	"github.com/spf13/cobra"
)

var userPath string

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&userPath, "path", "p", "", "Path to the user folder")
	importCmd.MarkFlagRequired("path")
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data from v1",
	Long:  `Import data from v1`,
	Run: func(cmd *cobra.Command, args []string) {
		db := infra.NewSqliteDB(DbPath)
		app := App(db)

		posts, err := importer.AllUserPosts(userPath)
		if err != nil {
			panic(err)
		}

		for _, post := range posts {
			fmt.Println(post.Meta.Type)
			switch post.Meta.Type {
			case "article":
				article := model.Article{}
				article.SetID(post.Id)
				article.SetMetaData(model.ArticleMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
				})
				article.SetPublishedAt(&post.Meta.Date)
				app.EntryService.Create(&article)

			case "bookmark":

			case "reply":

			case "photo":

			case "note":

			case "recipe":

			case "page":

			default:
				panic("Unknown type")
			}

		}
	},
}
