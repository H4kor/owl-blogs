package owl

import (
	"bytes"
	"io/ioutil"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"gopkg.in/yaml.v2"
)

type Post struct {
	user  *User
	id    string
	title string
}

type PostMeta struct {
	Title   string   `yaml:"title"`
	Aliases []string `yaml:"aliases"`
	Date    string   `yaml:"date"`
	Draft   bool     `yaml:"draft"`
}

func (post Post) Id() string {
	return post.id
}

func (post Post) Dir() string {
	return path.Join(post.user.Dir(), "public", post.id)
}

func (post Post) MediaDir() string {
	return path.Join(post.Dir(), "media")
}

func (post Post) UrlPath() string {
	return post.user.UrlPath() + "posts/" + post.id + "/"
}

func (post Post) FullUrl() string {
	return post.user.FullUrl() + "posts/" + post.id + "/"
}

func (post Post) UrlMediaPath(filename string) string {
	return post.UrlPath() + "media/" + filename
}

func (post Post) Title() string {
	return post.title
}

func (post Post) ContentFile() string {
	return path.Join(post.Dir(), "index.md")
}

func (post Post) Content() []byte {
	// read file
	data, _ := ioutil.ReadFile(post.ContentFile())
	return data
}

func (post Post) MarkdownData() (bytes.Buffer, PostMeta) {
	data := post.Content()

	// get yaml metadata block
	meta := PostMeta{}
	trimmedData := bytes.TrimSpace(data)
	// check first line is ---
	if string(trimmedData[0:4]) == "---\n" {
		trimmedData = trimmedData[4:]
		// find --- end
		end := bytes.Index(trimmedData, []byte("\n---\n"))
		if end != -1 {
			metaData := trimmedData[:end]
			yaml.Unmarshal(metaData, &meta)
			data = trimmedData[end+5:]
		}
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			// meta.Meta,
			extension.GFM,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(data, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	// metaData := meta.Get(context)

	return buf, meta

}

func (post Post) Aliases() []string {
	_, metaData := post.MarkdownData()
	return metaData.Aliases
}
