package kiss

import (
	"bytes"

	"github.com/yuin/goldmark"
)

func RenderPost(post Post) string {

	var buf bytes.Buffer
	if err := goldmark.Convert(post.Content(), &buf); err != nil {
		panic(err)
	}
	return buf.String()
}
