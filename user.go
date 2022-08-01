package kiss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type User struct {
	repo Repository
	name string
}

type UserConfig struct {
	Title       string `yaml:"title"`
	SubTitle    string `yaml:"subtitle"`
	HeaderColor string `yaml:"header_color"`
}

func (user User) Dir() string {
	return path.Join(user.repo.UsersDir(), user.name)
}

func (user User) Path() string {
	return "/user/" + user.name
}

func (user User) PostDir() string {
	return path.Join(user.Dir(), "public")
}

func (user User) MetaDir() string {
	return path.Join(user.Dir(), "meta")
}

func (user User) ConfigFile() string {
	return path.Join(user.MetaDir(), "config.yml")
}

func (user User) Name() string {
	return user.name
}

func (user User) Posts() ([]string, error) {
	postFiles := walkDir(path.Join(user.Dir(), "public"))
	posts := make([]string, 0)
	for _, id := range postFiles {
		if strings.HasSuffix(id, "/index.md") {
			posts = append(posts, id[:len(id)-9])
		}
	}
	return posts, nil
}

func (user User) GetPost(id string) (Post, error) {
	post := Post{user: user, id: id}
	_, metaData := post.MarkdownData()
	title := metaData["title"]
	post.title = fmt.Sprint(title)

	return post, nil
}

func (user User) CreateNewPost(title string) (Post, error) {
	timestamp := time.Now().UTC().Unix()
	folder_name := fmt.Sprintf("%d-%s", timestamp, title)
	post_dir := path.Join(user.Dir(), "public", folder_name)

	// if post already exists, add -n to the end of the name
	i := 0
	for {
		if dirExists(post_dir) {
			i++
			folder_name = fmt.Sprintf("%d-%s-%d", timestamp, title, i)
			post_dir = path.Join(user.Dir(), "public", folder_name)
		} else {
			break
		}
	}
	post := Post{user: user, id: folder_name, title: title}

	initial_content := ""
	initial_content += "---\n"
	initial_content += "title: " + title + "\n"
	initial_content += "---\n"
	initial_content += "\n"
	initial_content += "Write your post here.\n"

	// create post file
	os.Mkdir(post_dir, 0755)
	os.WriteFile(post.ContentFile(), []byte(initial_content), 0644)
	// create media dir
	os.Mkdir(post.MediaDir(), 0755)
	return post, nil
}

func (user User) Template() (string, error) {
	// load base.html
	path := path.Join(user.Dir(), "meta", "base.html")
	base_html, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(base_html), nil
}

func (user User) Config() (UserConfig, error) {
	config_path := user.ConfigFile()
	config_data, err := ioutil.ReadFile(config_path)
	if err != nil {
		return UserConfig{}, err
	}
	var meta UserConfig
	err = yaml.Unmarshal(config_data, &meta)
	if err != nil {
		return UserConfig{}, err
	}
	return meta, nil
}

func (user User) SetConfig(new_config UserConfig) error {
	config_path := user.ConfigFile()
	config_data, err := yaml.Marshal(new_config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config_path, config_data, 0644)
	if err != nil {
		return err
	}
	return nil
}
