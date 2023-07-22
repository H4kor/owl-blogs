package forms

import (
	"fmt"
	"mime/multipart"
	"owl-blogs/app"
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

type Form struct {
	data   interface{}
	binSvc *app.BinaryService
}

type FormFieldParams struct {
	InputType string
	Widget    string
}

type FormField struct {
	Name   string
	Value  string
	Params FormFieldParams
}

func NewForm(data interface{}, binaryService *app.BinaryService) *Form {
	return &Form{
		data:   data,
		binSvc: binaryService,
	}
}

func (s *FormFieldParams) ApplyTag(tagKey string, tagValue string) error {
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

func (s *FormField) Html() string {
	html := ""
	html += fmt.Sprintf("<label for=\"%v\">%v</label>\n", s.Name, s.Name)
	if s.Params.InputType == "text" && s.Params.Widget == "textarea" {
		html += fmt.Sprintf("<textarea name=\"%v\" id=\"%v\" rows=\"20\">%v</textarea>\n", s.Name, s.Name, s.Value)
	} else {
		html += fmt.Sprintf("<input type=\"%v\" name=\"%v\" id=\"%v\" value=\"%v\" />\n", s.Params.InputType, s.Name, s.Name, s.Value)
	}
	return html
}

func FieldToFormField(field reflect.StructField, value string) (FormField, error) {
	formField := FormField{
		Name:   field.Name,
		Value:  value,
		Params: FormFieldParams{},
	}
	tag := field.Tag.Get("owl")
	for _, param := range strings.Split(tag, " ") {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			continue
		}
		err := formField.Params.ApplyTag(parts[0], parts[1])
		if err != nil {
			return FormField{}, err
		}
	}
	return formField, nil
}

func StructToFormFields(data interface{}) ([]FormField, error) {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	dataType := reflect.TypeOf(data).Elem()
	numFields := dataType.NumField()
	fields := []FormField{}
	for i := 0; i < numFields; i++ {
		field, err := FieldToFormField(dataType.Field(i), dataValue.FieldByIndex([]int{i}).String())
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (s *Form) HtmlForm() (string, error) {
	fields, err := StructToFormFields(s.data)
	if err != nil {
		return "", err
	}

	html := "<form method=\"POST\" enctype=\"multipart/form-data\">\n"
	for _, field := range fields {
		html += field.Html()
	}
	html += "<input type=\"submit\" value=\"Submit\" />\n"
	html += "</form>\n"

	return html, nil
}

func (s *Form) Parse(ctx HttpFormData) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil context")
	}
	dataVal := reflect.ValueOf(s.data)
	if dataVal.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("meta data is not a pointer")
	}
	fields, err := StructToFormFields(s.data)
	if err != nil {
		return nil, err
	}
	for _, field := range fields {
		fieldName := field.Name

		if field.Params.InputType == "file" {
			file, err := ctx.FormFile(fieldName)
			if err != nil {
				return nil, err
			}
			fileData, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer fileData.Close()
			fileBytes := make([]byte, file.Size)
			_, err = fileData.Read(fileBytes)
			if err != nil {
				return nil, err
			}

			binaryFile, err := s.binSvc.Create(file.Filename, fileBytes)
			if err != nil {
				return nil, err
			}

			metaField := dataVal.Elem().FieldByName(fieldName)
			if metaField.IsValid() {
				metaField.SetString(binaryFile.Id)
			}
		} else {
			formValue := ctx.FormValue(fieldName)
			metaField := dataVal.Elem().FieldByName(fieldName)
			if metaField.IsValid() {
				metaField.SetString(formValue)
			}
		}

	}

	return s.data, nil
}
