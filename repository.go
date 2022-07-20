package kiss

import (
	"fmt"
	"os"
	"path"
	"time"
)

type Repository struct {
	name string
}

type User struct {
	repo Repository
	name string
}

type Post struct {
	user  User
	title string
}

func CreateRepository(name string) (Repository, error) {
	newRepo := Repository{name: name}
	// check if repository already exists
	if dirExists(newRepo.Dir()) {
		return Repository{}, fmt.Errorf("Repository already exists")
	}

	os.Mkdir(name, 0755)
	return newRepo, nil
}

func (repo Repository) Dir() string {
	return repo.name
}

func (user User) Dir() string {
	return path.Join(user.repo.Dir(), user.name)
}

func PostDir(post Post) string {
	return path.Join(post.user.Dir(), "public", post.title)
}

func CreateNewUser(repo Repository, name string) (User, error) {
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
	os.WriteFile(path.Join(user_dir, "meta", "VERSION"), []byte("0.0.1"), 0644)
	os.WriteFile(path.Join(user_dir, "meta", "base.html"), []byte("<html><body><{{content}}/body></html>"), 0644)

	return new_user, nil
}

func CreateNewPost(user User, title string) {
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

	initial_content := "# " + title
	// create post file
	os.Mkdir(post_dir, 0755)
	os.WriteFile(path.Join(post_dir, "index.md"), []byte(initial_content), 0644)
}
