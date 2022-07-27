package kiss

import (
	"bytes"
	"io/ioutil"
	"path"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	user  User
	id    string
	title string
}

func (post Post) Dir() string {
	return path.Join(post.user.Dir(), "public", post.id)
}

func (post Post) UrlPath() string {
	return post.user.Path() + "/" + post.id
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

func (post Post) MarkdownData() (bytes.Buffer, map[string]interface{}) {
	data := post.Content()
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			extension.GFM,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(data, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)

	return buf, metaData

}
