package kiss

import (
	"fmt"
	"os"
	"path"
)

type Repository struct {
	name string
}

func CreateRepository(name string) (Repository, error) {
	newRepo := Repository{name: name}
	// check if repository already exists
	if dirExists(newRepo.Dir()) {
		return Repository{}, fmt.Errorf("Repository already exists")
	}

	os.Mkdir(name, 0755)
	os.Mkdir(path.Join(name, "users"), 0755)
	return newRepo, nil
}

func OpenRepository(name string) (Repository, error) {

	repo := Repository{name: name}
	if !dirExists(repo.Dir()) {
		return Repository{}, fmt.Errorf("Repository does not exist")
	}

	return repo, nil

}

func (repo Repository) Dir() string {
	return repo.name
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

	base_template := "<html><body><{{content}}/body></html>"

	// create Meta files
	os.WriteFile(path.Join(user_dir, "meta", "VERSION"), []byte("0.0.1"), 0644)
	os.WriteFile(path.Join(user_dir, "meta", "base.html"), []byte(base_template), 0644)

	return new_user, nil
}
