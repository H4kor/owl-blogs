package forms

import (
	"fmt"
	"reflect"
	"strings"
)

type Widget interface {
	Html() string
	ParseValue(value string, output reflect.Value) error
}

type TextWidget struct {
	FormField
}

func (s *TextWidget) Html() string {
	html := ""
	html += fmt.Sprintf("<input type=\"text\" name=\"%v\" value=\"%v\">\n", s.Name, s.Value.String())
	return html
}

func (s *TextWidget) ParseValue(value string, output reflect.Value) error {
	output.SetString(value)
	return nil
}

type TextareaWidget struct {
	FormField
}

func (s *TextareaWidget) Html() string {
	html := ""
	html += fmt.Sprintf("<textarea name=\"%v\" rows=\"20\">%v</textarea>\n", s.Name, s.Value.String())
	return html
}

func (s *TextareaWidget) ParseValue(value string, output reflect.Value) error {
	output.SetString(value)
	return nil
}

type TextListWidget struct {
	FormField
}

func (s *TextListWidget) Html() string {
	valueList := s.Value.Interface().([]string)
	value := strings.Join(valueList, "\n")

	html := ""
	html += fmt.Sprintf("<textarea name=\"%v\" rows=\"20\">%v</textarea>\n", s.Name, value)
	return html
}

func (s *TextListWidget) ParseValue(value string, output reflect.Value) error {
	list := strings.Split(value, "\n")
	// trim entries
	for i, item := range list {
		list[i] = strings.TrimSpace(item)
	}
	// remove empty entries
	for i := len(list) - 1; i >= 0; i-- {
		if list[i] == "" {
			list = append(list[:i], list[i+1:]...)
		}
	}

	output.Set(reflect.ValueOf(list))
	return nil
}
