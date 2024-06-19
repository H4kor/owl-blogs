package render

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"net/url"
	"owl-blogs/domain/model"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/hashtag"
)

type TemplateData struct {
	Data       interface{}
	SiteConfig model.SiteConfig
}

//go:embed templates
var templates embed.FS
var SiteConfigService model.SiteConfigInterface

var funcMap = template.FuncMap{
	"markdown": func(text string) template.HTML {
		html, err := RenderMarkdown(text)
		if err != nil {
			return ">>>could not render markdown<<<"
		}
		return template.HTML(html)
	},
	"urljoin": func(elems ...string) string {
		r, _ := url.JoinPath(elems[0], elems[1:]...)
		return r
	},
}

func CreateTemplateWithBase(templateName string) (*template.Template, error) {

	return template.New(templateName).Funcs(funcMap).ParseFS(
		templates,
		"templates/base.tmpl",
		"templates/"+templateName+".tmpl",
	)
}

func RenderTemplateWithBase(w io.Writer, templateName string, data interface{}) error {

	t, err := CreateTemplateWithBase(templateName)

	if err != nil {
		return err
	}

	siteConfig, err := SiteConfigService.GetSiteConfig()
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, "base", TemplateData{
		Data:       data,
		SiteConfig: siteConfig,
	})

	return err

}

func RenderTemplateToString(templateName string, data interface{}) (template.HTML, error) {
	tmplStr, _ := templates.ReadFile("templates/" + templateName + ".tmpl")

	t, err := template.New("templates/" + templateName + ".tmpl").Funcs(funcMap).Parse(string(tmplStr))

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = t.Execute(&output, data)
	return template.HTML(output.String()), err
}

type HashTagResolver struct {
}

// ResolveHashtag reports the link that the provided hashtag Node
// should point to, or an empty destination for hashtags that should
// not link to anything.
func (*HashTagResolver) ResolveHashtag(node *hashtag.Node) (destination []byte, err error) {
	return []byte("/tags/" + string(node.Tag) + "/"), nil
}

func RenderMarkdown(mdText string) (string, error) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			&hashtag.Extender{
				Resolver: &HashTagResolver{},
				Attributes: []hashtag.Attribute{
					{Attr: "class", Value: "p-category"},
				},
			},
			&ThumbnailExtension{},
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	err := markdown.Convert([]byte(mdText), &buf, parser.WithContext(context))

	return buf.String(), err

}
