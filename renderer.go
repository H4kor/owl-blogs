package owl

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"
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
	Post    Post
	Content template.HTML
}

type AuthRequestData struct {
	Me                  string
	ClientId            string
	RedirectUri         string
	State               string
	Scope               string
	Scopes              []string // Split version of scope. filled by rendering function.
	ResponseType        string
	CodeChallenge       string
	CodeChallengeMethod string
	User                User
	CsrfToken           string
}

type EditorViewData struct {
	User      User
	Error     string
	CsrfToken string
}

type ErrorMessage struct {
	Error   string
	Message string
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func listUrl(user User, id string) string {
	return user.ListUrl(PostList{
		Id: id,
	})
}

func postUrl(user User, id string) string {
	post, _ := user.GetPost(id)
	return post.UrlPath()
}

func renderEmbedTemplate(templateFile string, data interface{}) (string, error) {
	templateStr, err := embed_files.ReadFile(templateFile)
	if err != nil {
		return "", err
	}
	return renderTemplateStr(templateStr, data)
}

func renderTemplateStr(templateStr []byte, data interface{}) (string, error) {
	t, err := template.New("_").Funcs(template.FuncMap{
		"noescape": noescape,
		"listUrl":  listUrl,
		"postUrl":  postUrl,
	}).Parse(string(templateStr))
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
	t, err := template.New("index").Funcs(template.FuncMap{
		"noescape": noescape,
		"listUrl":  listUrl,
		"postUrl":  postUrl,
	}).Parse(baseTemplate)
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
	err = t.Execute(&html, full_data)
	return html.String(), err
}

func renderPostContent(post Post) (string, error) {
	buf := post.RenderedContent()
	postHtml, err := renderEmbedTemplate(
		fmt.Sprintf("embed/%s/detail.html", post.TemplateDir()),
		PostRenderData{
			Title:   post.Title(),
			Post:    post,
			Content: template.HTML(buf),
		},
	)
	return postHtml, err
}

func RenderPost(post Post) (string, error) {
	postHtml, err := renderPostContent(post)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(*post.User(), PageContent{
		Title:       post.Title(),
		Description: post.Meta().Description,
		Content:     template.HTML(postHtml),
		Type:        "article",
		SelfUrl:     post.FullUrl(),
	})
}

func RenderIndexPage(user User) (string, error) {
	posts, _ := user.PrimaryFeedPosts()

	postHtml, err := renderEmbedTemplate("embed/post-list.html", posts)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(user, PageContent{
		Title:   "Index",
		Content: template.HTML(postHtml),
	})
}

func RenderPostList(user User, list *PostList) (string, error) {
	posts, _ := user.GetPostsOfList(*list)
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
	reqData.Scopes = strings.Split(reqData.Scope, " ")
	authHtml, err := renderEmbedTemplate("embed/auth.html", reqData)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(reqData.User, PageContent{
		Title:   "Auth",
		Content: template.HTML(authHtml),
	})
}

func RenderUserError(user User, error ErrorMessage) (string, error) {
	errHtml, err := renderEmbedTemplate("embed/error.html", error)
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(user, PageContent{
		Title:   "Error",
		Content: template.HTML(errHtml),
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

func RenderLoginPage(user User, error_type string, csrfToken string) (string, error) {
	loginHtml, err := renderEmbedTemplate("embed/editor/login.html", EditorViewData{
		User:      user,
		Error:     error_type,
		CsrfToken: csrfToken,
	})
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(user, PageContent{
		Title:   "Login",
		Content: template.HTML(loginHtml),
	})
}

func RenderEditorPage(user User, csrfToken string) (string, error) {
	editorHtml, err := renderEmbedTemplate("embed/editor/editor.html", EditorViewData{
		User:      user,
		CsrfToken: csrfToken,
	})
	if err != nil {
		return "", err
	}

	return renderIntoBaseTemplate(user, PageContent{
		Title:   "Editor",
		Content: template.HTML(editorHtml),
	})
}
