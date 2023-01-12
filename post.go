package owl

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"path"
	"sort"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v2"
)

type GenericPost struct {
	user       *User
	id         string
	metaLoaded bool
	meta       PostMeta
	wmLock     sync.Mutex
}

func (post *GenericPost) TemplateDir() string {
	return post.Meta().Type
}

type Post interface {
	TemplateDir() string

	// Actual Data
	User() *User
	Id() string
	Title() string
	Meta() PostMeta
	Content() []byte
	RenderedContent() string
	Aliases() []string

	// Filesystem
	Dir() string
	MediaDir() string
	ContentFile() string

	// Urls
	UrlPath() string
	FullUrl() string
	UrlMediaPath(filename string) string

	// Webmentions Support
	IncomingWebmentions() []WebmentionIn
	OutgoingWebmentions() []WebmentionOut
	PersistIncomingWebmention(webmention WebmentionIn) error
	PersistOutgoingWebmention(webmention *WebmentionOut) error
	AddIncomingWebmention(source string) error
	EnrichWebmention(webmention WebmentionIn) error
	ApprovedIncomingWebmentions() []WebmentionIn
	ScanForLinks() error
	SendWebmention(webmention WebmentionOut) error
}

type ReplyData struct {
	Url  string `yaml:"url"`
	Text string `yaml:"text"`
}
type BookmarkData struct {
	Url  string `yaml:"url"`
	Text string `yaml:"text"`
}

type RecipeData struct {
	Yield       string   `yaml:"yield"`
	Duration    string   `yaml:"duration"`
	Ingredients []string `yaml:"ingredients"`
}

type PostMeta struct {
	Type        string       `yaml:"type"`
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	Aliases     []string     `yaml:"aliases"`
	Date        time.Time    `yaml:"date"`
	Draft       bool         `yaml:"draft"`
	Reply       ReplyData    `yaml:"reply"`
	Bookmark    BookmarkData `yaml:"bookmark"`
	Recipe      RecipeData   `yaml:"recipe"`
}

func (pm PostMeta) FormattedDate() string {
	return pm.Date.Format("02-01-2006 15:04:05")
}

func (pm *PostMeta) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type T struct {
		Type        string       `yaml:"type"`
		Title       string       `yaml:"title"`
		Description string       `yaml:"description"`
		Aliases     []string     `yaml:"aliases"`
		Draft       bool         `yaml:"draft"`
		Reply       ReplyData    `yaml:"reply"`
		Bookmark    BookmarkData `yaml:"bookmark"`
		Recipe      RecipeData   `yaml:"recipe"`
	}
	type S struct {
		Date string `yaml:"date"`
	}

	var t T
	var s S
	if err := unmarshal(&t); err != nil {
		return err
	}
	if err := unmarshal(&s); err != nil {
		return err
	}

	pm.Type = t.Type
	if pm.Type == "" {
		pm.Type = "article"
	}
	pm.Title = t.Title
	pm.Description = t.Description
	pm.Aliases = t.Aliases
	pm.Draft = t.Draft
	pm.Reply = t.Reply
	pm.Bookmark = t.Bookmark
	pm.Recipe = t.Recipe

	possibleFormats := []string{
		"2006-01-02",
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}

	for _, format := range possibleFormats {
		if t, err := time.Parse(format, s.Date); err == nil {
			pm.Date = t
			break
		}
	}

	return nil
}

type PostWebmetions struct {
	Incoming []WebmentionIn  `ymal:"incoming"`
	Outgoing []WebmentionOut `ymal:"outgoing"`
}

func (post *GenericPost) Id() string {
	return post.id
}

func (post *GenericPost) User() *User {
	return post.user
}

func (post *GenericPost) Dir() string {
	return path.Join(post.user.Dir(), "public", post.id)
}

func (post *GenericPost) IncomingWebmentionsFile() string {
	return path.Join(post.Dir(), "incoming_webmentions.yml")
}

func (post *GenericPost) OutgoingWebmentionsFile() string {
	return path.Join(post.Dir(), "outgoing_webmentions.yml")
}

func (post *GenericPost) MediaDir() string {
	return path.Join(post.Dir(), "media")
}

func (post *GenericPost) UrlPath() string {
	return post.user.UrlPath() + "posts/" + post.id + "/"
}

func (post *GenericPost) FullUrl() string {
	return post.user.FullUrl() + "posts/" + post.id + "/"
}

func (post *GenericPost) UrlMediaPath(filename string) string {
	return post.UrlPath() + "media/" + filename
}

func (post *GenericPost) Title() string {
	return post.Meta().Title
}

func (post *GenericPost) ContentFile() string {
	return path.Join(post.Dir(), "index.md")
}

func (post *GenericPost) Meta() PostMeta {
	if !post.metaLoaded {
		post.LoadMeta()
	}
	return post.meta
}

func (post *GenericPost) Content() []byte {
	// read file
	data, _ := os.ReadFile(post.ContentFile())
	return data
}

func (post *GenericPost) RenderedContent() string {
	data := post.Content()

	// trim yaml block
	// TODO this can be done nicer
	trimmedData := bytes.TrimSpace(data)
	// ensure that data ends with a newline
	trimmedData = append(trimmedData, []byte("\n")...)
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
	if config, _ := post.user.repo.Config(); config.AllowRawHtml {
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

	return buf.String()

}

func (post *GenericPost) Aliases() []string {
	return post.Meta().Aliases
}

func (post *GenericPost) LoadMeta() error {
	data := post.Content()

	// get yaml metadata block
	meta := PostMeta{}
	trimmedData := bytes.TrimSpace(data)
	// ensure that data ends with a newline
	trimmedData = append(trimmedData, []byte("\n")...)
	// check first line is ---
	if string(trimmedData[0:4]) == "---\n" {
		trimmedData = trimmedData[4:]
		// find --- end
		end := bytes.Index(trimmedData, []byte("---\n"))
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

func (post *GenericPost) IncomingWebmentions() []WebmentionIn {
	// return parsed webmentions
	fileName := post.IncomingWebmentionsFile()
	if !fileExists(fileName) {
		return []WebmentionIn{}
	}

	webmentions := []WebmentionIn{}
	loadFromYaml(fileName, &webmentions)

	return webmentions
}

func (post *GenericPost) OutgoingWebmentions() []WebmentionOut {
	// return parsed webmentions
	fileName := post.OutgoingWebmentionsFile()
	if !fileExists(fileName) {
		return []WebmentionOut{}
	}

	webmentions := []WebmentionOut{}
	loadFromYaml(fileName, &webmentions)

	return webmentions
}

// PersistWebmentionOutgoing persists incoming webmention
func (post *GenericPost) PersistIncomingWebmention(webmention WebmentionIn) error {
	post.wmLock.Lock()
	defer post.wmLock.Unlock()

	wms := post.IncomingWebmentions()

	// if target is not in status, add it
	replaced := false
	for i, t := range wms {
		if t.Source == webmention.Source {
			wms[i].UpdateWith(webmention)
			replaced = true
			break
		}
	}

	if !replaced {
		wms = append(wms, webmention)
	}

	err := saveToYaml(post.IncomingWebmentionsFile(), wms)
	if err != nil {
		return err
	}

	return nil
}

// PersistOutgoingWebmention persists a webmention to the webmention file.
func (post *GenericPost) PersistOutgoingWebmention(webmention *WebmentionOut) error {
	post.wmLock.Lock()
	defer post.wmLock.Unlock()

	wms := post.OutgoingWebmentions()

	// if target is not in webmention, add it
	replaced := false
	for i, t := range wms {
		if t.Target == webmention.Target {
			wms[i].UpdateWith(*webmention)
			replaced = true
			break
		}
	}

	if !replaced {
		wms = append(wms, *webmention)
	}

	err := saveToYaml(post.OutgoingWebmentionsFile(), wms)
	if err != nil {
		return err
	}

	return nil
}

func (post *GenericPost) AddIncomingWebmention(source string) error {
	// Check if file already exists
	wm := WebmentionIn{
		Source: source,
	}

	defer func() {
		go post.EnrichWebmention(wm)
	}()
	return post.PersistIncomingWebmention(wm)
}

func (post *GenericPost) EnrichWebmention(webmention WebmentionIn) error {
	resp, err := post.user.repo.HttpClient.Get(webmention.Source)
	if err == nil {
		entry, err := post.user.repo.Parser.ParseHEntry(resp)
		if err == nil {
			webmention.Title = entry.Title
			return post.PersistIncomingWebmention(webmention)
		}
	}
	return err
}

func (post *GenericPost) ApprovedIncomingWebmentions() []WebmentionIn {
	webmentions := post.IncomingWebmentions()
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

// ScanForLinks scans the post content for links and adds them to the
// `status.yml` file for the post. The links are not scanned by this function.
func (post *GenericPost) ScanForLinks() error {
	// this could be done in markdown parsing, but I don't want to
	// rely on goldmark for this (yet)
	postHtml := post.RenderedContent()
	links, _ := post.user.repo.Parser.ParseLinksFromString(postHtml)
	// add reply url if set
	if post.Meta().Reply.Url != "" {
		links = append(links, post.Meta().Reply.Url)
	}
	for _, link := range links {
		post.PersistOutgoingWebmention(&WebmentionOut{
			Target: link,
		})
	}
	return nil
}

func (post *GenericPost) SendWebmention(webmention WebmentionOut) error {
	defer post.PersistOutgoingWebmention(&webmention)

	// if last scan is less than 7 days ago, don't send webmention
	if webmention.ScannedAt.After(time.Now().Add(-7*24*time.Hour)) && !webmention.Supported {
		return errors.New("did not scan. Last scan was less than 7 days ago")
	}

	webmention.ScannedAt = time.Now()

	resp, err := post.user.repo.HttpClient.Get(webmention.Target)
	if err != nil {
		webmention.Supported = false
		return err
	}

	endpoint, err := post.user.repo.Parser.GetWebmentionEndpoint(resp)
	if err != nil {
		webmention.Supported = false
		return err
	}
	webmention.Supported = true

	// send webmention
	payload := url.Values{}
	payload.Set("source", post.FullUrl())
	payload.Set("target", webmention.Target)
	_, err = post.user.repo.HttpClient.PostForm(endpoint, payload)

	if err != nil {
		return err
	}

	// update webmention status
	webmention.LastSentAt = time.Now()
	return nil
}
