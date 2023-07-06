package editor_test

import (
	"mime/multipart"
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

type MockFormData struct {
}

func (f *MockFormData) FormFile(key string) (*multipart.FileHeader, error) {
	return nil, nil
}

func (f *MockFormData) FormValue(key string, defaultValue ...string) string {
	return key
}

type MockEntry struct {
	id          string
	publishedAt *time.Time
	metaData    MockEntryMetaData
}

func (e *MockEntry) ID() string {
	return e.id
}

func (e *MockEntry) Content() model.EntryContent {
	return model.EntryContent(e.metaData.Content)
}

func (e *MockEntry) PublishedAt() *time.Time {
	return e.publishedAt
}

func (e *MockEntry) MetaData() interface{} {
	return &e.metaData
}

func (e *MockEntry) Create(id string, publishedAt *time.Time, metaData model.EntryMetaData) error {
	e.id = id
	e.publishedAt = publishedAt
	e.metaData = *metaData.(*MockEntryMetaData)
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
	form := editor.NewEntryForm(&MockEntry{})
	html, err := form.HtmlForm()
	require.NoError(t, err)
	require.Contains(t, html, "<form")
	require.Contains(t, html, "method=\"POST\"")
	require.Contains(t, html, "<input type=\"file\" name=\"Image\"")
	require.Contains(t, html, "<input type=\"text\" name=\"Content\"")
	require.Contains(t, html, "<input type=\"submit\" value=\"Submit\"")

}

func TestFormParseNil(t *testing.T) {
	form := editor.NewEntryForm(&MockEntry{})
	_, err := form.Parse(nil)
	require.Error(t, err)
}

func TestFormParse(t *testing.T) {
	form := editor.NewEntryForm(&MockEntry{})
	entry, err := form.Parse(&MockFormData{})
	require.NoError(t, err)
	require.Equal(t, "Image", entry.MetaData().(*MockEntryMetaData).Image)
	require.Equal(t, "Content", entry.MetaData().(*MockEntryMetaData).Content)
}
