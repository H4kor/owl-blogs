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
	Title string
	Post  template.HTML
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

	user_config, _ := user.Config()
	full_data := struct {
		Title        string
		Content      template.HTML
		UserTitle    string
		UserSubtitle string
		HeaderColor  string
	}{
		Title:        data.Title,
		Content:      data.Content,
		UserTitle:    user_config.Title,
		UserSubtitle: user_config.SubTitle,
		HeaderColor:  user_config.HeaderColor,
	}

	var html bytes.Buffer
	t.Execute(&html, full_data)

	return html.String(), nil
}

func RenderPost(post Post) (string, error) {
	buf, _ := post.MarkdownData()
	postHtml, err := renderEmbedTemplate("embed/post.html", PostRenderData{
		Title: post.Title(),
		Post:  template.HTML(buf.String()),
	})
	if err != nil {
		return "", err
	}

	data := PageContent{
		Title:   post.Title(),
		Content: template.HTML(postHtml),
	}

	return renderIntoBaseTemplate(*post.user, data)
}

func RenderIndexPage(user User) (string, error) {
	posts, _ := user.Posts()

	postHtml := ""
	for _, post := range posts {
		postHtml += "<h2><a href=\"" + post.UrlPath() + "\">" + post.Title() + "</a></h2>\n"
	}

	data := PageContent{
		Title:   "Index",
		Content: template.HTML(postHtml),
	}

	return renderIntoBaseTemplate(user, data)

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
