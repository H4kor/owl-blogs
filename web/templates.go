package web

import (
	"embed"
	"io"
	"text/template"
)

//go:embed templates
var templates embed.FS

func CreateTemplate(templateName string) (*template.Template, error) {

	return template.ParseFS(
		templates,
		"templates/base.tmpl",
		"templates/"+templateName+".tmpl",
	)

}

func RenderTemplate(w io.Writer, templateName string, data interface{}) error {

	t, err := CreateTemplate(templateName)

	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, "base", data)

	return err

}
