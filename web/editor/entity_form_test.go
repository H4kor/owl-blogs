package editor_test

import (
	"owl-blogs/domain/model"
	"owl-blogs/web/editor"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockEntryMetaData struct {
	Image   string `owl:"inputType=file"`
	Content string `owl:"inputType=text"`
}

type MockEntry struct {
	id          string
	content     model.EntryContent
	publishedAt *time.Time
	metaData    *MockEntryMetaData
}

func (e *MockEntry) ID() string {
	return e.id
}

func (e *MockEntry) Content() model.EntryContent {
	return e.content
}

func (e *MockEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *MockEntry) MetaData() interface{} {
	return e.metaData
}

func (e *MockEntry) Create(id string, content string, publishedAt *time.Time, metaData model.EntryMetaData) error {
	e.id = id
	e.content = model.EntryContent(content)
	e.publishedAt = publishedAt
	e.metaData = metaData.(*MockEntryMetaData)
	return nil
}

func TestFieldToFormField(t *testing.T) {
	field := reflect.TypeOf(&MockEntryMetaData{}).Elem().Field(0)
	formField, err := editor.FieldToFormField(field)
	require.NoError(t, err)
	require.Equal(t, "Image", formField.Name)
	require.Equal(t, "file", formField.Params.InputType)
}

func TestStructToFields(t *testing.T) {
	fields, err := editor.StructToFormFields(&MockEntryMetaData{})
	require.NoError(t, err)
	require.Len(t, fields, 2)
	require.Equal(t, "Image", fields[0].Name)
	require.Equal(t, "file", fields[0].Params.InputType)
	require.Equal(t, "Content", fields[1].Name)
	require.Equal(t, "text", fields[1].Params.InputType)
}

func TestEditorEntryForm_HtmlForm(t *testing.T) {
	formService := editor.NewEditorFormService(&MockEntry{})
	form, err := formService.HtmlForm()
	require.NoError(t, err)
	require.Contains(t, form, "<form")
	require.Contains(t, form, "method=\"POST\"")
	require.Contains(t, form, "<input type=\"file\" name=\"Image\" />")
	require.Contains(t, form, "<input type=\"text\" name=\"Content\" />")
	require.Contains(t, form, "<input type=\"submit\" value=\"Submit\" />")

}
