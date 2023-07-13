package render

import (
	"bytes"
	"embed"
	"io"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed templates
var templates embed.FS

var funcMap = template.FuncMap{
	"markdown": func(text string) string {
		html, err := RenderMarkdown(text)
		if err != nil {
			return ">>>could not render markdown<<<"
		}
		return html
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

	err = t.ExecuteTemplate(w, "base", data)

	return err

}

func RenderTemplateToString(templateName string, data interface{}) (string, error) {
	tmplStr, _ := templates.ReadFile("templates/" + templateName + ".tmpl")

	t, err := template.New("templates/" + templateName + ".tmpl").Funcs(funcMap).Parse(string(tmplStr))

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = t.Execute(&output, data)
	return output.String(), err
}

func RenderMarkdown(mdText string) (string, error) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			// meta.Meta,
			extension.GFM,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	err := markdown.Convert([]byte(mdText), &buf, parser.WithContext(context))

	return buf.String(), err

}
