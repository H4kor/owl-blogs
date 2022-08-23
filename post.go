package owl

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type Post struct {
	user       *User
	id         string
	title      string
	metaLoaded bool
	meta       PostMeta
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

func (post Post) WebmentionDir() string {
	return path.Join(post.Dir(), "webmention")
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

func (post *Post) Meta() PostMeta {
	if !post.metaLoaded {
		post.LoadMeta()
	}
	return post.meta
}

func (post Post) Content() []byte {
	// read file
	data, _ := ioutil.ReadFile(post.ContentFile())
	return data
}

func (post Post) RenderedContent() bytes.Buffer {
	data := post.Content()

	// trim yaml block
	// TODO this can be done nicer
	trimmedData := bytes.TrimSpace(data)
	// check first line is ---
	if string(trimmedData[0:4]) == "---\n" {
		trimmedData = trimmedData[4:]
		// find --- end
		end := bytes.Index(trimmedData, []byte("\n---\n"))
		if end != -1 {
			data = trimmedData[end+5:]
		}
	}

	options := goldmark.WithRendererOptions()
	if post.user.repo.AllowRawHtml() {
		options = goldmark.WithRendererOptions(
			html.WithUnsafe(),
		)
	}

	markdown := goldmark.New(
		options,
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

	return buf

}

func (post Post) Aliases() []string {
	return post.Meta().Aliases
}

func (post *Post) LoadMeta() error {
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
			err := yaml.Unmarshal(metaData, &meta)
			if err != nil {
				return err
			}
		}
	}

	post.meta = meta
	return nil
}

func (post *Post) AddWebmention(source string) error {
	hash := sha256.Sum256([]byte(source))
	hashStr := base64.URLEncoding.EncodeToString(hash[:])
	data := "source: " + source
	return os.WriteFile(path.Join(post.WebmentionDir(), hashStr+".yml"), []byte(data), 0644)
}

func (post *Post) Webmentions() []string {
	return listDir(post.WebmentionDir())
}
