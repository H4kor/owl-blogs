package kiss

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed embed/user-list.html
var userListTemplateStr string

//go:embed embed/post.html
var postTemplateStr string

type PageContent struct {
	Title   string
	Content template.HTML
}

type PostRenderData struct {
	Title string
	Post  template.HTML
}

func RenderPost(post Post) (string, error) {
	baseTemplate, _ := post.user.Template()
	buf, _ := post.MarkdownData()

	postTemplate, _ := template.New("post").Parse(postTemplateStr)
	var postHtml bytes.Buffer
	err := postTemplate.Execute(&postHtml, PostRenderData{
		Title: post.Title(),
		Post:  template.HTML(buf.String()),
	})
	if err != nil {
		return "", err
	}

	data := PageContent{
		Title:   post.Title(),
		Content: template.HTML(postHtml.String()),
	}

	var html bytes.Buffer
	t, err := template.New("page").Parse(baseTemplate)

	t.Execute(&html, data)

	return html.String(), err
}

func RenderIndexPage(user User) (string, error) {
	baseTemplate, _ := user.Template()
	posts, _ := user.Posts()

	postHtml := ""
	for _, postId := range posts {
		post, _ := user.GetPost(postId)
		postHtml += "<h2><a href=\"" + post.Path() + "\">" + post.Title() + "</a></h2>\n"
	}

	data := PageContent{
		Title:   "Index",
		Content: template.HTML(postHtml),
	}

	var html bytes.Buffer
	t, err := template.New("post").Parse(baseTemplate)

	t.Execute(&html, data)

	return html.String(), err

}

// func RenderUserList(user User) (string, error) {
// 	base_template, _ := user.Template()
// 	users, _ := user.repo.Users()
// 	template.New("user_list").Parse()
// 	return strings.Replace(template, "{{content}}", userHtml, -1), nil
// }
