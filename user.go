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

func (user User) Dir() string {
	return path.Join(user.repo.Dir(), "users", user.name)
}

func (user User) Name() string {
	return user.name
}

func (user User) Posts() ([]Post, error) {
	postNames := listDir(path.Join(user.Dir(), "public"))
	posts := make([]Post, len(postNames))
	for i, name := range postNames {
		posts[i] = Post{user: user, id: name}
	}
	return posts, nil
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
	post := Post{user: user, id: folder_name}

	initial_content := "# " + title
	// create post file
	os.Mkdir(post_dir, 0755)
	os.WriteFile(post.ContentFile(), []byte(initial_content), 0644)
	return post, nil
}
