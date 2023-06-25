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

type EntryFormField struct {
	Name   string
	Params map[string]string
}

func NewEditorFormService(entry model.Entry) *EditorEntryForm {
	return &EditorEntryForm{
		entry: entry,
	}
}

func (s *EditorEntryForm) HtmlForm() string {
	meta := s.entry.MetaData()
	entryType := reflect.TypeOf(meta).Elem()
	numFields := entryType.NumField()

	fields := []EntryFormField{}
	for i := 0; i < numFields; i++ {
		field := EntryFormField{
			Name:   entryType.Field(i).Name,
			Params: map[string]string{},
		}
		tag := entryType.Field(i).Tag.Get("owl")
		for _, param := range strings.Split(tag, " ") {
			parts := strings.Split(param, "=")
			if len(parts) == 2 {
				field.Params[parts[0]] = parts[1]
			} else {
				field.Params[param] = ""
			}
		}
		fields = append(fields, field)
	}

	return fmt.Sprintf("%v", fields)
}
