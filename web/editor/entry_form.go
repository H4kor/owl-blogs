package editor

import (
	"fmt"
	"mime/multipart"
	"owl-blogs/domain/model"
	"reflect"
	"strings"
)

type HttpFormData interface {
	// FormFile returns the first file by key from a MultipartForm.
	FormFile(key string) (*multipart.FileHeader, error)
	// FormValue returns the first value by key from a MultipartForm.
	// Search is performed in QueryArgs, PostArgs, MultipartForm and FormFile in this particular order.
	// Defaults to the empty string "" if the form value doesn't exist.
	// If a default value is given, it will return that value if the form value does not exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	FormValue(key string, defaultValue ...string) string
}

type EditorEntryForm struct {
	entry model.Entry
}

type EntryFormFieldParams struct {
	InputType string
	Widget    string
}

type EntryFormField struct {
	Name   string
	Params EntryFormFieldParams
}

func NewEntryForm(entry model.Entry) *EditorEntryForm {
	return &EditorEntryForm{
		entry: entry,
	}
}

func (s *EntryFormFieldParams) ApplyTag(tagKey string, tagValue string) error {
	switch tagKey {
	case "inputType":
		s.InputType = tagValue
	case "widget":
		s.Widget = tagValue
	default:
		return fmt.Errorf("unknown tag key: %v", tagKey)
	}
	return nil
}

func (s *EntryFormField) Html() string {
	html := ""
	html += fmt.Sprintf("<label for=\"%v\">%v</label>\n", s.Name, s.Name)
	if s.Params.InputType == "text" && s.Params.Widget == "textarea" {
		html += fmt.Sprintf("<textarea name=\"%v\" id=\"%v\" rows=\"20\"></textarea>\n", s.Name, s.Name)
	} else {
		html += fmt.Sprintf("<input type=\"%v\" name=\"%v\" id=\"%v\" />\n", s.Params.InputType, s.Name, s.Name)
	}
	return html
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

func (s *EditorEntryForm) Parse(ctx HttpFormData) (model.Entry, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	meta := s.entry.MetaData()
	metaVal := reflect.ValueOf(meta)
	if metaVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("meta data is not a pointer")
	}
	fields, err := StructToFormFields(meta)
	if err != nil {
		return nil, err
	}
	for field := range fields {
		fieldName := fields[field].Name
		fieldValue := ctx.FormValue(fieldName)
		metaField := metaVal.Elem().FieldByName(fieldName)
		if metaField.IsValid() {
			metaField.SetString(fieldValue)
		}
	}

	return s.entry, nil
}
