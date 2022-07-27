package kiss

import (
	"embed"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

//go:embed embed/initial/base.html
var base_template string

//go:embed embed/*
var static_files embed.FS

var VERSION = "0.0.1"

type Repository struct {
	name string
}

func CreateRepository(name string) (Repository, error) {
	newRepo := Repository{name: name}
	// check if repository already exists
	if dirExists(newRepo.Dir()) {
		return Repository{}, fmt.Errorf("Repository already exists")
	}

	os.Mkdir(newRepo.Dir(), 0755)
	os.Mkdir(newRepo.UsersDir(), 0755)
	os.Mkdir(newRepo.StaticDir(), 0755)

	// copy all files from static_files embed.FS to StaticDir
	staticFiles, _ := static_files.ReadDir("embed/initial/static")
	for _, file := range staticFiles {
		if file.IsDir() {
			continue
		}
		src_data, _ := static_files.ReadFile(file.Name())
		os.WriteFile(newRepo.StaticDir()+"/"+file.Name(), src_data, 0644)
	}

	// copy repo_base.html to base.html
	src_data, _ := static_files.ReadFile("embed/initial/repo_base.html")
	os.WriteFile(newRepo.Dir()+"/base.html", src_data, 0644)
	return newRepo, nil
}

func OpenRepository(name string) (Repository, error) {

	repo := Repository{name: name}
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

func (repo Repository) Template() (string, error) {
	// load base.html
	path := path.Join(repo.Dir(), "base.html")
	base_html, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(base_html), nil
}

func (repo Repository) Users() ([]User, error) {
	userNames := listDir(repo.UsersDir())
	users := make([]User, len(userNames))
	for i, name := range userNames {
		users[i] = User{repo: repo, name: name}
	}
	return users, nil
}

func (repo Repository) CreateUser(name string) (User, error) {
	new_user := User{repo: repo, name: name}
	// check if user already exists
	if dirExists(new_user.Dir()) {
		return User{}, fmt.Errorf("User already exists")
	}

	// creates repo/name folder if it doesn't exist
	user_dir := new_user.Dir()
	os.Mkdir(user_dir, 0755)
	os.Mkdir(path.Join(user_dir, "meta"), 0755)
	// create public folder
	os.Mkdir(path.Join(user_dir, "public"), 0755)

	// create Meta files
	os.WriteFile(path.Join(user_dir, "meta", "VERSION"), []byte(VERSION), 0644)
	os.WriteFile(path.Join(user_dir, "meta", "base.html"), []byte(base_template), 0644)

	meta, _ := yaml.Marshal(UserConfig{
		Title:       name,
		SubTitle:    "",
		HeaderColor: "#bdd6be",
	})
	os.WriteFile(new_user.ConfigFile(), meta, 0644)

	return new_user, nil
}

func (repo Repository) GetUser(name string) (User, error) {
	user := User{repo: repo, name: name}
	if !dirExists(user.Dir()) {
		return User{}, fmt.Errorf("User does not exist")
	}
	return user, nil
}
