package importer

import (
	"bytes"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"
)

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
	PhotoPath   string       `yaml:"photo"`
}

type Post struct {
	Id      string
	Meta    PostMeta
	Content string
}

func (post *Post) MediaDir() string {
	return path.Join("public", post.Id, "media")
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
		PhotoPath   string       `yaml:"photo"`
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
	pm.PhotoPath = t.PhotoPath

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

func LoadContent(data []byte) string {

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

	return string(data)
}

func LoadMeta(data []byte) (PostMeta, error) {

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
				return PostMeta{}, err
			}
		}
	}

	return meta, nil
}

func AllUserPosts(userPath string) ([]Post, error) {
	postFiles := ListDir(path.Join(userPath, "public"))
	posts := make([]Post, 0)
	for _, id := range postFiles {
		// if is a directory and has index.md, add to posts
		if dirExists(path.Join(userPath, "public", id)) {
			if fileExists(path.Join(userPath, "public", id, "index.md")) {
				postData, err := os.ReadFile(path.Join(userPath, "public", id, "index.md"))
				if err != nil {
					return nil, err
				}
				meta, err := LoadMeta(postData)
				if err != nil {
					return nil, err
				}
				post := Post{
					Id:      id,
					Content: LoadContent(postData),
					Meta:    meta,
				}
				posts = append(posts, post)
			}
		}
	}

	return posts, nil
}

func ListDir(path string) []string {
	dir, _ := os.Open(path)
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)
	return files
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
