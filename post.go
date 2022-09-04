package owl

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"

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

type PostStatus struct {
	Webmentions []WebmentionOut
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

func (post *Post) WebmentionFile(source string) string {

	hash := sha256.Sum256([]byte(source))
	hashStr := base64.URLEncoding.EncodeToString(hash[:])
	return path.Join(post.WebmentionDir(), hashStr+".yml")
}

func (post *Post) PersistWebmention(webmention WebmentionIn) error {
	// ensure dir exists
	os.MkdirAll(post.WebmentionDir(), 0755)

	// write to file
	fileName := post.WebmentionFile(webmention.Source)
	data, err := yaml.Marshal(webmention)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

func (post *Post) Webmention(source string) (WebmentionIn, error) {
	// ensure dir exists
	os.MkdirAll(post.WebmentionDir(), 0755)

	// Check if file exists
	fileName := post.WebmentionFile(source)
	if !fileExists(fileName) {
		// return error if file doesn't exist
		return WebmentionIn{}, fmt.Errorf("Webmention file not found: %s", source)
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return WebmentionIn{}, err
	}

	mention := WebmentionIn{}
	err = yaml.Unmarshal(data, &mention)
	if err != nil {
		return WebmentionIn{}, err
	}

	return mention, nil
}

func (post *Post) AddWebmention(source string) error {
	// Check if file already exists
	_, err := post.Webmention(source)
	if err != nil {
		webmention := WebmentionIn{
			Source: source,
		}
		defer post.EnrichWebmention(source)
		return post.PersistWebmention(webmention)
	}
	return nil
}

func (post *Post) EnrichWebmention(source string) error {
	html, err := post.user.repo.Retriever.Get(source)
	if err == nil {
		webmention, err := post.Webmention(source)
		if err != nil {
			return err
		}
		entry, err := post.user.repo.Parser.ParseHEntry(html)
		if err == nil {
			webmention.Title = entry.Title
			return post.PersistWebmention(webmention)
		}
	}
	return err
}

func (post *Post) Webmentions() []WebmentionIn {
	// ensure dir exists
	os.MkdirAll(post.WebmentionDir(), 0755)
	files := listDir(post.WebmentionDir())
	webmentions := []WebmentionIn{}
	for _, file := range files {
		data, err := os.ReadFile(path.Join(post.WebmentionDir(), file))
		if err != nil {
			continue
		}
		mention := WebmentionIn{}
		err = yaml.Unmarshal(data, &mention)
		if err != nil {
			continue
		}
		webmentions = append(webmentions, mention)
	}

	return webmentions
}

func (post *Post) ApprovedWebmentions() []WebmentionIn {
	webmentions := post.Webmentions()
	approved := []WebmentionIn{}
	for _, webmention := range webmentions {
		if webmention.ApprovalStatus == "approved" {
			approved = append(approved, webmention)
		}
	}

	// sort by retrieved date
	sort.Slice(approved, func(i, j int) bool {
		return approved[i].RetrievedAt.After(approved[j].RetrievedAt)
	})
	return approved
}
