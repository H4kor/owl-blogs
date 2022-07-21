package kiss

import "strings"

func RenderPost(post Post) (string, error) {
	template, _ := post.user.Template()
	buf, _ := post.MarkdownData()
	postHtml := "<h1>" + post.Title() + "</h1>\n"
	postHtml += buf.String()
	return strings.Replace(template, "{{content}}", postHtml, -1), nil
}
