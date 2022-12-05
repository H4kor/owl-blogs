package owl

import (
	_ "embed"
	"fmt"
	"os"
	"path"
)

//go:embed embed/initial/base.html
var base_template string

var VERSION = "0.0.1"

type Repository struct {
	name       string
	HttpClient HttpClient
	Parser     HtmlParser
}

type RepoConfig struct {
	Domain       string `yaml:"domain"`
	SingleUser   string `yaml:"single_user"`
	AllowRawHtml bool   `yaml:"allow_raw_html"`
}

func CreateRepository(name string, config RepoConfig) (Repository, error) {
	newRepo := Repository{name: name, Parser: OwlHtmlParser{}, HttpClient: &OwlHttpClient{}}
	// check if repository already exists
	if dirExists(newRepo.Dir()) {
		return Repository{}, fmt.Errorf("Repository already exists")
	}

	os.Mkdir(newRepo.Dir(), 0755)
	os.Mkdir(newRepo.UsersDir(), 0755)
	os.Mkdir(newRepo.StaticDir(), 0755)

	// create config file
	if config.Domain == "" {
		config.Domain = "http://localhost:8080"
	}
	saveToYaml(path.Join(newRepo.Dir(), "config.yml"), config)

	// copy all files from static_files embed.FS to StaticDir
	staticFiles, _ := embed_files.ReadDir("embed/initial/static")
	for _, file := range staticFiles {
		if file.IsDir() {
			continue
		}
		src_data, _ := embed_files.ReadFile("embed/initial/static/" + file.Name())
		os.WriteFile(newRepo.StaticDir()+"/"+file.Name(), src_data, 0644)
	}

	// copy repo/ to newRepo.Dir()
	init_files, _ := embed_files.ReadDir("embed/initial/repo")
	for _, file := range init_files {
		if file.IsDir() {
			continue
		}
		src_data, _ := embed_files.ReadFile("embed/initial/repo/" + file.Name())
		os.WriteFile(newRepo.Dir()+"/"+file.Name(), src_data, 0644)
	}
	return newRepo, nil
}

func OpenRepository(name string) (Repository, error) {

	repo := Repository{name: name, Parser: OwlHtmlParser{}, HttpClient: &OwlHttpClient{}}
	if !dirExists(repo.Dir()) {
		return Repository{}, fmt.Errorf("Repository does not exist: " + repo.Dir())
	}

	return repo, nil
}

func (repo Repository) Dir() string {
	return repo.name
}

func (repo Repository) StaticDir() string {
	return path.Join(repo.Dir(), "static")
}

func (repo Repository) UsersDir() string {
	return path.Join(repo.Dir(), "users")
}

func (repo Repository) UserUrlPath(user User) string {
	config, _ := repo.Config()
	if config.SingleUser != "" {
		return "/"
	}
	return "/user/" + user.name + "/"
}

func (repo Repository) FullUrl() string {
	config, _ := repo.Config()
	return config.Domain
}

func (repo Repository) Template() (string, error) {
	// load base.html
	path := path.Join(repo.Dir(), "base.html")
	base_html, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(base_html), nil
}

func (repo Repository) Users() ([]User, error) {
	config, _ := repo.Config()
	if config.SingleUser != "" {
		return []User{{repo: &repo, name: config.SingleUser}}, nil
	}

	userNames := listDir(repo.UsersDir())
	users := make([]User, len(userNames))
	for i, name := range userNames {
		users[i] = User{repo: &repo, name: name}
	}
	return users, nil
}

func (repo *Repository) CreateUser(name string) (User, error) {
	new_user := User{repo: repo, name: name}
	// check if user already exists
	if dirExists(new_user.Dir()) {
		return User{}, fmt.Errorf("User already exists")
	}

	// creates repo/name folder if it doesn't exist
	user_dir := new_user.Dir()
	os.Mkdir(user_dir, 0755)
	// create folders
	os.Mkdir(path.Join(user_dir, "meta"), 0755)
	os.Mkdir(path.Join(user_dir, "public"), 0755)
	os.Mkdir(path.Join(user_dir, "media"), 0755)

	// create Meta files
	os.WriteFile(path.Join(user_dir, "meta", "VERSION"), []byte(VERSION), 0644)
	os.WriteFile(path.Join(user_dir, "meta", "base.html"), []byte(base_template), 0644)

	saveToYaml(new_user.ConfigFile(), UserConfig{
		Title:       name,
		SubTitle:    "",
		HeaderColor: "#bdd6be",
	})

	return new_user, nil
}

func (repo Repository) GetUser(name string) (User, error) {
	user := User{repo: &repo, name: name}
	if !dirExists(user.Dir()) {
		return User{}, fmt.Errorf("User does not exist")
	}
	return user, nil
}

func (repo Repository) PostAliases() (map[string]IPost, error) {
	users, err := repo.Users()
	if err != nil {
		return nil, err
	}
	aliases := make(map[string]IPost)
	for _, user := range users {
		user_aliases, err := user.PostAliases()
		if err != nil {
			return nil, err
		}
		for alias, post := range user_aliases {
			aliases[alias] = post
		}
	}
	return aliases, nil
}

func (repo Repository) Config() (RepoConfig, error) {
	meta := RepoConfig{}
	err := loadFromYaml(path.Join(repo.Dir(), "config.yml"), &meta)
	return meta, err

}
