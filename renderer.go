package owl

import (
	"bytes"
	_ "embed"
	"html/template"
)

type PageContent struct {
	Title       string
	Description string
	Content     template.HTML
	Type        string
	SelfUrl     string
}

type PostRenderData struct {
	Title   string
	Post    *Post
	Content template.HTML
}

type AuthRequestData struct {
	Me           string
	ClientId     string
	RedirectUri  string
	State        string
	ResponseType string
	User         User
	CsrfToken    string
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
		Title       string
		Description string
		Content     template.HTML
		Type        string
		SelfUrl     string
		User        User
	}{
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		Type:        data.Type,
		SelfUrl:     data.SelfUrl,
		User:        user,
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
		Title:       post.Title(),
		Description: post.Meta().Description,
		Content:     template.HTML(postHtml),
		Type:        "article",
		SelfUrl:     post.FullUrl(),
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

func RenderUserAuthPage(reqData AuthRequestData) (string, error) {
	authHtml, err := renderEmbedTemplate("embed/auth.html", reqData)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(reqData.User, PageContent{
		Title:   "Auth",
		Content: template.HTML(authHtml),
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
