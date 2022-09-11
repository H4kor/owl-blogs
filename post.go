package owl

import (
	"bytes"
	"errors"
	"io/ioutil"
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

type Post struct {
	user       *User
	id         string
	title      string
	metaLoaded bool
	meta       PostMeta
	wmLock     sync.Mutex
}

type PostMeta struct {
	Title   string    `yaml:"title"`
	Aliases []string  `yaml:"aliases"`
	Date    time.Time `yaml:"date"`
	Draft   bool      `yaml:"draft"`
}

func (pm PostMeta) FormattedDate() string {
	return pm.Date.Format("02-01-2006 15:04:05")
}

func (pm *PostMeta) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type T struct {
		Title   string   `yaml:"title"`
		Aliases []string `yaml:"aliases"`
		Draft   bool     `yaml:"draft"`
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

	pm.Title = t.Title
	pm.Aliases = t.Aliases
	pm.Draft = t.Draft

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

func (post *Post) Id() string {
	return post.id
}

func (post *Post) Dir() string {
	return path.Join(post.user.Dir(), "public", post.id)
}

func (post *Post) WebmentionsFile() string {
	return path.Join(post.Dir(), "webmentions.yml")
}

func (post *Post) MediaDir() string {
	return path.Join(post.Dir(), "media")
}

func (post *Post) UrlPath() string {
	return post.user.UrlPath() + "posts/" + post.id + "/"
}

func (post *Post) FullUrl() string {
	return post.user.FullUrl() + "posts/" + post.id + "/"
}

func (post *Post) UrlMediaPath(filename string) string {
	return post.UrlPath() + "media/" + filename
}

func (post *Post) Title() string {
	return post.title
}

func (post *Post) ContentFile() string {
	return path.Join(post.Dir(), "index.md")
}

func (post *Post) Meta() PostMeta {
	if !post.metaLoaded {
		post.LoadMeta()
	}
	return post.meta
}

func (post *Post) Content() []byte {
	// read file
	data, _ := ioutil.ReadFile(post.ContentFile())
	return data
}

func (post *Post) RenderedContent() bytes.Buffer {
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

	return buf

}

func (post *Post) Aliases() []string {
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

// Webmentions returns list of incoming and outgoing webmentions
func (post *Post) Webmentions() PostWebmetions {
	// read status file
	// return parsed webmentions
	fileName := post.WebmentionsFile()
	if !fileExists(fileName) {
		return PostWebmetions{}
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return PostWebmetions{}
	}

	webmentions := PostWebmetions{}
	err = yaml.Unmarshal(data, &webmentions)
	if err != nil {
		return PostWebmetions{}
	}

	return webmentions
}

func (post *Post) IncomingWebmentions() []WebmentionIn {
	return post.Webmentions().Incoming
}

func (post *Post) OutgoingWebmentions() []WebmentionOut {
	return post.Webmentions().Outgoing
}

func (post *Post) persistWebmentions(webmentions PostWebmetions) error {
	data, err := yaml.Marshal(webmentions)
	if err != nil {
		return err
	}

	err = os.WriteFile(post.WebmentionsFile(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// PersistWebmentionOutgoing persists incoming webmention
func (post *Post) PersistIncomingWebmention(webmention WebmentionIn) error {
	post.wmLock.Lock()
	defer post.wmLock.Unlock()

	wms := post.Webmentions()

	// if target is not in status, add it
	replaced := false
	for i, t := range wms.Incoming {
		if t.Source == webmention.Source {
			wms.Incoming[i].UpdateWith(webmention)
			replaced = true
			break
		}
	}

	if !replaced {
		wms.Incoming = append(wms.Incoming, webmention)
	}

	return post.persistWebmentions(wms)
}

// PersistOutgoingWebmention persists a webmention to the webmention file.
func (post *Post) PersistOutgoingWebmention(webmention *WebmentionOut) error {
	post.wmLock.Lock()
	defer post.wmLock.Unlock()

	wms := post.Webmentions()

	// if target is not in webmention, add it
	replaced := false
	for i, t := range wms.Outgoing {
		if t.Target == webmention.Target {
			wms.Outgoing[i].UpdateWith(*webmention)
			replaced = true
			break
		}
	}

	if !replaced {
		wms.Outgoing = append(wms.Outgoing, *webmention)
	}

	return post.persistWebmentions(wms)
}

func (post *Post) AddIncomingWebmention(source string) error {
	// Check if file already exists
	wm := WebmentionIn{
		Source: source,
	}

	defer func() {
		go post.EnrichWebmention(wm)
	}()
	return post.PersistIncomingWebmention(wm)
}

func (post *Post) EnrichWebmention(webmention WebmentionIn) error {
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

func (post *Post) ApprovedIncomingWebmentions() []WebmentionIn {
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
func (post *Post) ScanForLinks() error {
	// this could be done in markdown parsing, but I don't want to
	// rely on goldmark for this (yet)
	postHtml := post.RenderedContent()
	links, _ := post.user.repo.Parser.ParseLinksFromString(postHtml.String())
	for _, link := range links {
		post.PersistOutgoingWebmention(&WebmentionOut{
			Target: link,
		})
	}
	return nil
}

func (post *Post) SendWebmention(webmention WebmentionOut) error {
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
