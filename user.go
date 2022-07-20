package kiss

import (
	"fmt"
	"os"
	"path"
	"time"
)

type User struct {
	repo Repository
	name string
}

type Post struct {
	user  User
	title string
}

func (user User) Dir() string {
	return path.Join(user.repo.Dir(), "users", user.name)
}

func (user User) Name() string {
	return user.name
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

func PostDir(post Post) string {
	return path.Join(post.user.Dir(), "public", post.title)
}
