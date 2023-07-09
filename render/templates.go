package render

import (
	"bytes"
	"embed"
	"io"
	"text/template"
)

//go:embed templates
var templates embed.FS

func CreateTemplateWithBase(templateName string) (*template.Template, error) {

	return template.ParseFS(
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

	t, err := template.ParseFS(
		templates,
		"templates/"+templateName+".tmpl",
	)

	if err != nil {
		return "", err
	}

	var output bytes.Buffer

	err = t.Execute(&output, data)
	return output.String(), err
}
