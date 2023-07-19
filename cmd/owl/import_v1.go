package main

import (
	"fmt"
	"os"
	"owl-blogs/config"
	"owl-blogs/domain/model"
	entrytypes "owl-blogs/entry_types"
	"owl-blogs/importer"
	"owl-blogs/infra"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

		// import config
		bytes, err := os.ReadFile(path.Join(userPath, "meta/config.yml"))
		if err != nil {
			panic(err)
		}
		v1Config := importer.V1UserConfig{}
		yaml.Unmarshal(bytes, &v1Config)

		mes := []model.MeLinks{}
		for _, me := range v1Config.Me {
			mes = append(mes, model.MeLinks{
				Name: me.Name,
				Url:  me.Url,
			})
		}

		lists := []model.EntryList{}
		for _, list := range v1Config.Lists {
			lists = append(lists, model.EntryList{
				Id:       list.Id,
				Title:    list.Title,
				Include:  importer.ConvertTypeList(list.Include, app.Registry),
				ListType: list.ListType,
			})
		}

		headerMenu := []model.MenuItem{}
		for _, item := range v1Config.HeaderMenu {
			headerMenu = append(headerMenu, model.MenuItem{
				Title: item.Title,
				List:  item.List,
				Url:   item.Url,
				Post:  item.Post,
			})
		}

		footerMenu := []model.MenuItem{}
		for _, item := range v1Config.FooterMenu {
			footerMenu = append(footerMenu, model.MenuItem{
				Title: item.Title,
				List:  item.List,
				Url:   item.Url,
				Post:  item.Post,
			})
		}

		v2Config := &model.SiteConfig{}
		err = app.SiteConfigRepo.Get(config.SITE_CONFIG, v2Config)
		if err != nil {
			panic(err)
		}
		v2Config.Title = v1Config.Title
		v2Config.SubTitle = v1Config.SubTitle
		v2Config.HeaderColor = v1Config.HeaderColor
		v2Config.AuthorName = v1Config.AuthorName
		v2Config.Me = mes
		v2Config.Lists = lists
		v2Config.PrimaryListInclude = importer.ConvertTypeList(v1Config.PrimaryListInclude, app.Registry)
		v2Config.HeaderMenu = headerMenu
		v2Config.FooterMenu = footerMenu

		err = app.SiteConfigRepo.Update(config.SITE_CONFIG, v2Config)
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
				entry = &entrytypes.Bookmark{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.BookmarkMetaData{
					Url:     post.Meta.Bookmark.Url,
					Title:   post.Meta.Bookmark.Text,
					Content: post.Content,
				})
			case "reply":
				entry = &entrytypes.Reply{}
				entry.SetID(post.Id)
				entry.SetPublishedAt(&post.Meta.Date)
				entry.SetMetaData(&entrytypes.ReplyMetaData{
					Url:     post.Meta.Reply.Url,
					Title:   post.Meta.Reply.Text,
					Content: post.Content,
				})
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
