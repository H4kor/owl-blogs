package kiss

import (
	"io/ioutil"
	"path"
)

type Post struct {
	user User
	id   string
}

func (post Post) Dir() string {
	return path.Join(post.user.Dir(), "public", post.id)
}

func (post Post) ContentFile() string {
	return path.Join(post.Dir(), "index.md")
}

func (post Post) Content() []byte {
	// read file
	data, _ := ioutil.ReadFile(post.ContentFile())
	return data
}
