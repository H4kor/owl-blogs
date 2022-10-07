package owl

import (
	"bytes"
	_ "embed"
	"html/template"
)

type PageContent struct {
	Title   string
	Content template.HTML
}

type PostRenderData struct {
	Title   string
	Post    *Post
	Content template.HTML
}

func renderEmbedTemplate(templateFile string, data interface{}) (string, error) {
	templateStr, err := embed_files.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	return renderTemplateStr(templateStr, data)
}

func renderTemplateStr(templateStr []byte, data interface{}) (string, error) {
	t, err := template.New("_").Parse(string(templateStr))
	if err != nil {
		return "", err
	}
	var html bytes.Buffer
	err = t.Execute(&html, data)
	if err != nil {
		return "", err
	}
	return html.String(), nil
}

func renderIntoBaseTemplate(user User, data PageContent) (string, error) {
	baseTemplate, _ := user.Template()
	t, err := template.New("index").Parse(baseTemplate)
	if err != nil {
		return "", err
	}

	full_data := struct {
		Title        string
		Content      template.HTML
		User         User
		UserSubtitle string
		HeaderColor  string
	}{
		Title:   data.Title,
		Content: data.Content,
		User:    user,
	}

	var html bytes.Buffer
	t.Execute(&html, full_data)

	return html.String(), nil
}

func renderPostContent(post *Post) (string, error) {
	buf := post.RenderedContent()
	postHtml, err := renderEmbedTemplate("embed/post.html", PostRenderData{
		Title:   post.Title(),
		Post:    post,
		Content: template.HTML(buf.String()),
	})
	return postHtml, err
}

func RenderPost(post *Post) (string, error) {
	postHtml, err := renderPostContent(post)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(*post.user, PageContent{
		Title:   post.Title(),
		Content: template.HTML(postHtml),
	})
}

func RenderIndexPage(user User) (string, error) {
	posts, _ := user.Posts()

	postHtml, err := renderEmbedTemplate("embed/post-list.html", posts)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(user, PageContent{
		Title:   "Index",
		Content: template.HTML(postHtml),
	})

}

func RenderUserList(repo Repository) (string, error) {
	baseTemplate, _ := repo.Template()
	users, _ := repo.Users()
	userHtml, err := renderEmbedTemplate("embed/user-list.html", users)
	if err != nil {
		return "", err
	}

	data := PageContent{
		Title:   "Index",
		Content: template.HTML(userHtml),
	}

	return renderTemplateStr([]byte(baseTemplate), data)
}
