package main

import (
	"fmt"
	"os"
	"owl-blogs/domain/model"
	"owl-blogs/importer"
	"owl-blogs/infra"
	"path"

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
				entry := &model.Article{}
				entry.SetID(post.Id)

				fileData, err := os.ReadFile(path.Join(mediaDir, file))
				if err != nil {
					panic(err)
				}
				app.BinaryService.CreateEntryFile(file, fileData, entry)
			}

			switch post.Meta.Type {
			case "article":
				article := model.Article{}
				article.SetID(post.Id)
				article.SetPublishedAt(&post.Meta.Date)
				article.SetMetaData(&model.ArticleMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
				})
				app.EntryService.Create(&article)
			case "bookmark":

			case "reply":

			case "photo":
				photo := model.Image{}
				photo.SetID(post.Id)
				photo.SetPublishedAt(&post.Meta.Date)
				photo.SetMetaData(&model.ImageMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
					ImageId: post.Meta.PhotoPath,
				})
				app.EntryService.Create(&photo)
			case "note":
				note := model.Note{}
				note.SetID(post.Id)
				note.SetPublishedAt(&post.Meta.Date)
				note.SetMetaData(&model.NoteMetaData{
					Content: post.Content,
				})
				app.EntryService.Create(&note)
			case "recipe":
				recipe := model.Recipe{}
				recipe.SetID(post.Id)
				recipe.SetPublishedAt(&post.Meta.Date)
				recipe.SetMetaData(&model.RecipeMetaData{
					Title:       post.Meta.Title,
					Yield:       post.Meta.Recipe.Yield,
					Duration:    post.Meta.Recipe.Duration,
					Ingredients: post.Meta.Recipe.Ingredients,
					Content:     post.Content,
				})
				app.EntryService.Create(&recipe)
			case "page":
				page := model.Page{}
				page.SetID(post.Id)
				page.SetPublishedAt(&post.Meta.Date)
				page.SetMetaData(&model.PageMetaData{
					Title:   post.Meta.Title,
					Content: post.Content,
				})
				app.EntryService.Create(&page)
			default:
				panic("Unknown type")
			}

		}
	},
}
