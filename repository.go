package kiss

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path"
)

//go:embed embed/initial/base.html
var base_template string

//go:embed embed/initial/static/*
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

func (repo Repository) Users() ([]User, error) {
	userNames := listDir(path.Join(repo.Dir(), "users"))
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

	return new_user, nil
}

func (repo Repository) GetUser(name string) (User, error) {
	user := User{repo: repo, name: name}
	if !dirExists(user.Dir()) {
		return User{}, fmt.Errorf("User does not exist")
	}
	return user, nil
}
