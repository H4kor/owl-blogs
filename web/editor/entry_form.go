package editor

import (
	"fmt"
	"owl-blogs/domain/model"
	"reflect"
	"strings"
)

type EditorEntryForm struct {
	entry model.Entry
}

type EntryFormFieldParams struct {
	InputType string
}

type EntryFormField struct {
	Name   string
	Params EntryFormFieldParams
}

func NewEditorFormService(entry model.Entry) *EditorEntryForm {
	return &EditorEntryForm{
		entry: entry,
	}
}

func (s *EntryFormFieldParams) ApplyTag(tagKey string, tagValue string) error {
	switch tagKey {
	case "inputType":
		s.InputType = tagValue
	default:
		return fmt.Errorf("unknown tag key: %v", tagKey)
	}
	return nil
}

func (s *EntryFormField) Html() string {
	return fmt.Sprintf("<input type=\"%v\" name=\"%v\" />\n", s.Params.InputType, s.Name)
}

func FieldToFormField(field reflect.StructField) (EntryFormField, error) {
	formField := EntryFormField{
		Name:   field.Name,
		Params: EntryFormFieldParams{},
	}
	tag := field.Tag.Get("owl")
	for _, param := range strings.Split(tag, " ") {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			continue
		}
		err := formField.Params.ApplyTag(parts[0], parts[1])
		if err != nil {
			return EntryFormField{}, err
		}
	}
	return formField, nil
}

func StructToFormFields(meta interface{}) ([]EntryFormField, error) {
	entryType := reflect.TypeOf(meta).Elem()
	numFields := entryType.NumField()
	fields := []EntryFormField{}
	for i := 0; i < numFields; i++ {
		field, err := FieldToFormField(entryType.Field(i))
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (s *EditorEntryForm) HtmlForm() (string, error) {
	meta := s.entry.MetaData()
	fields, err := StructToFormFields(meta)
	if err != nil {
		return "", err
	}

	html := "<form method=\"POST\">\n"
	for _, field := range fields {
		html += field.Html()
	}
	html += "<input type=\"submit\" value=\"Submit\" />\n"
	html += "</form>\n"

	return html, nil
}
