package owl

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

type User struct {
	repo *Repository
	name string
}

type UserConfig struct {
	Title         string `yaml:"title"`
	SubTitle      string `yaml:"subtitle"`
	HeaderColor   string `yaml:"header_color"`
	TwitterHandle string `yaml:"twitter_handle"`
	GitHubHandle  string `yaml:"github_handle"`
	AuthorName    string `yaml:"author_name"`
}

func (user User) Dir() string {
	return path.Join(user.repo.UsersDir(), user.name)
}

func (user User) UrlPath() string {
	return user.repo.UserUrlPath(user)
}

func (user User) FullUrl() string {
	url, _ := url.JoinPath(user.repo.FullUrl(), user.UrlPath())
	return url
}

func (user User) WebmentionUrl() string {
	url, _ := url.JoinPath(user.FullUrl(), "webmention/")
	return url
}

func (user User) MediaUrl() string {
	url, _ := url.JoinPath(user.UrlPath(), "media")
	return url
}

func (user User) PostDir() string {
	return path.Join(user.Dir(), "public")
}

func (user User) MetaDir() string {
	return path.Join(user.Dir(), "meta")
}

func (user User) MediaDir() string {
	return path.Join(user.Dir(), "media")
}

func (user User) ConfigFile() string {
	return path.Join(user.MetaDir(), "config.yml")
}

func (user User) Name() string {
	return user.name
}

func (user User) AvatarUrl() string {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif"} {
		if fileExists(path.Join(user.MediaDir(), "avatar"+ext)) {
			url, _ := url.JoinPath(user.MediaUrl(), "avatar"+ext)
			return url
		}
	}
	return ""
}

func (user User) FaviconUrl() string {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif", ".ico"} {
		if fileExists(path.Join(user.MediaDir(), "favicon"+ext)) {
			url, _ := url.JoinPath(user.MediaUrl(), "favicon"+ext)
			return url
		}
	}
	return ""
}

func (user User) Posts() ([]*Post, error) {
	postFiles := listDir(path.Join(user.Dir(), "public"))
	posts := make([]*Post, 0)
	for _, id := range postFiles {
		// if is a directory and has index.md, add to posts
		if dirExists(path.Join(user.Dir(), "public", id)) {
			if fileExists(path.Join(user.Dir(), "public", id, "index.md")) {
				post, _ := user.GetPost(id)
				posts = append(posts, post)
			}
		}
	}

	// remove drafts
	n := 0
	for _, post := range posts {
		meta := post.Meta()
		if !meta.Draft {
			posts[n] = post
			n++
		}
	}
	posts = posts[:n]

	type PostWithDate struct {
		post *Post
		date time.Time
	}

	postDates := make([]PostWithDate, len(posts))
	for i, post := range posts {
		meta := post.Meta()
		postDates[i] = PostWithDate{post: post, date: meta.Date}
	}

	// sort posts by date
	sort.Slice(postDates, func(i, j int) bool {
		return postDates[i].date.After(postDates[j].date)
	})

	for i, post := range postDates {
		posts[i] = post.post
	}

	return posts, nil
}

func (user User) GetPost(id string) (*Post, error) {
	// check if posts index.md exists
	if !fileExists(path.Join(user.Dir(), "public", id, "index.md")) {
		return &Post{}, fmt.Errorf("post %s does not exist", id)
	}

	post := Post{user: &user, id: id}
	// post.loadMeta()
	meta := post.Meta()
	title := meta.Title
	post.title = fmt.Sprint(title)

	return &post, nil
}

func (user User) CreateNewPost(title string, draft bool) (*Post, error) {
	folder_name := toDirectoryName(title)
	post_dir := path.Join(user.Dir(), "public", folder_name)

	// if post already exists, add -n to the end of the name
	i := 0
	for {
		if dirExists(post_dir) {
			i++
			folder_name = toDirectoryName(fmt.Sprintf("%s-%d", title, i))
			post_dir = path.Join(user.Dir(), "public", folder_name)
		} else {
			break
		}
	}
	post := Post{user: &user, id: folder_name, title: title}
	meta := PostMeta{
		Title:   title,
		Date:    time.Now(),
		Aliases: []string{},
		Draft:   draft,
	}

	initial_content := ""
	initial_content += "---\n"
	// write meta
	meta_bytes, err := yaml.Marshal(meta)
	if err != nil {
		return &Post{}, err
	}
	initial_content += string(meta_bytes)
	initial_content += "---\n"
	initial_content += "\n"
	initial_content += "Write your post here.\n"

	// create post file
	os.Mkdir(post_dir, 0755)
	os.WriteFile(post.ContentFile(), []byte(initial_content), 0644)
	// create media dir
	os.Mkdir(post.MediaDir(), 0755)
	return &post, nil
}

func (user User) Template() (string, error) {
	// load base.html
	path := path.Join(user.Dir(), "meta", "base.html")
	base_html, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(base_html), nil
}

func (user User) Config() (UserConfig, error) {
	meta := UserConfig{}
	err := loadFromYaml(user.ConfigFile(), &meta)
	return meta, err
}

func (user User) SetConfig(new_config UserConfig) error {
	return saveToYaml(user.ConfigFile(), new_config)
}

func (user User) PostAliases() (map[string]*Post, error) {
	post_aliases := make(map[string]*Post)
	posts, err := user.Posts()
	if err != nil {
		return post_aliases, err
	}
	for _, post := range posts {
		if err != nil {
			return post_aliases, err
		}
		for _, alias := range post.Aliases() {
			post_aliases[alias] = post
		}
	}
	return post_aliases, nil
}
