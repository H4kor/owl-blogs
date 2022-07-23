package kiss

import "strings"

func RenderPost(post Post) (string, error) {
	template, _ := post.user.Template()
	buf, _ := post.MarkdownData()
	postHtml := "<h1>" + post.Title() + "</h1>\n"
	postHtml += buf.String()
	return strings.Replace(template, "{{content}}", postHtml, -1), nil
}

func RenderIndexPage(user User) (string, error) {
	template, _ := user.Template()
	posts, _ := user.Posts()
	postHtml := ""
	for _, postId := range posts {
		post, _ := user.GetPost(postId)
		postHtml += "<h2><a href=\"" + post.Path() + "\">" + post.Title() + "</a></h2>\n"
	}
	return strings.Replace(template, "{{content}}", postHtml, -1), nil
}
