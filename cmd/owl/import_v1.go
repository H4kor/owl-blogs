package main

import (
	"fmt"
	"os"
	"owl-blogs/domain/model"
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/importer"
	"owl-blogs/infra"
	"path"

	"github.com/spf13/cobra"
)

var userPath string
var author string

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVarP(&userPath, "path", "p", "", "Path to the user folder")
	importCmd.MarkFlagRequired("path")
	importCmd.Flags().StringVarP(&author, "author", "a", "", "The author name")
	importCmd.MarkFlagRequired("author")
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
			existing, _ := app.EntryService.FindById(post.Id)
			if existing != nil {
				continue
			}
			fmt.Println(post.Meta.Type)

			// import assets
			mediaDir := path.Join(userPath, post.MediaDir())
			println(mediaDir)
			files := importer.ListDir(mediaDir)
			for _, file := range files {
				// mock entry to pass to binary service
				entry := &entrytypes.Article{}
				entry.SetID(post.Id)

				fileData, err := os.ReadFile(path.Join(mediaDir, file))
				if err != nil {
					panic(err)
				}
				app.BinaryService.CreateEntryFile(file, fileData, entry)
			}

			var entry model.Entry

			switch post.Meta.Type {
			case "article":
				entry = &entrytypes.Article{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.ArticleMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
				})
			case "bookmark":

			case "reply":

			case "photo":
				entry = &entrytypes.Image{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.ImageMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
					ImageId: post.Meta.PhotoPath,
				})
			case "note":
				entry = &entrytypes.Note{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.NoteMetaData{
					Content: post.Content,
				})
			case "recipe":
				entry = &entrytypes.Recipe{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.RecipeMetaData{
					Title:       post.Meta.Title,
					Yield:       post.Meta.Recipe.Yield,
					Duration:    post.Meta.Recipe.Duration,
					Ingredients: post.Meta.Recipe.Ingredients,
					Content:     post.Content,
				})
			case "page":
				entry = &entrytypes.Page{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.PageMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
				})
			default:
				panic("Unknown type")
			}

			if entry != nil {
				entry.SetAuthorId(author)
				app.EntryService.Create(entry)
			}
		}
	},
}
