package editor_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/test"
	"owl-blogs/web/editor"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockEntryMetaData struct {
	Image   string `owl:"inputType=file"`
	Content string `owl:"inputType=text"`
}

type MockFormData struct {
	fileHeader *multipart.FileHeader
}

func NewMockFormData() *MockFormData {
	fileDir, _ := os.Getwd()
	fileName := "../../test/fixtures/test.png"
	filePath := path.Join(fileDir, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("ImagePath", filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}
	io.Copy(part, file)
	writer.Close()

	multipartForm := multipart.NewReader(body, writer.Boundary())
	formData, err := multipartForm.ReadForm(0)
	if err != nil {
		panic(err)
	}
	fileHeader := formData.File["ImagePath"][0]

	return &MockFormData{fileHeader: fileHeader}
}

func (f *MockFormData) FormFile(key string) (*multipart.FileHeader, error) {
	return f.fileHeader, nil
}

func (f *MockFormData) FormValue(key string, defaultValue ...string) string {
	return key
}

type MockEntry struct {
	model.EntryBase
	metaData MockEntryMetaData
}

func (e *MockEntry) Content() model.EntryContent {
	return model.EntryContent(e.metaData.Content)
}

func (e *MockEntry) MetaData() interface{} {
	return &e.metaData
}

func (e *MockEntry) SetMetaData(metaData interface{}) {
	e.metaData = *metaData.(*MockEntryMetaData)
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
	form := editor.NewEntryForm(&MockEntry{}, nil)
	html, err := form.HtmlForm()
	require.NoError(t, err)
	require.Contains(t, html, "<form")
	require.Contains(t, html, "method=\"POST\"")
	require.Contains(t, html, "<input type=\"file\" name=\"Image\"")
	require.Contains(t, html, "<input type=\"text\" name=\"Content\"")
	require.Contains(t, html, "<input type=\"submit\" value=\"Submit\"")

}

func TestFormParseNil(t *testing.T) {
	form := editor.NewEntryForm(&MockEntry{}, nil)
	_, err := form.Parse(nil)
	require.Error(t, err)
}

func TestFormParse(t *testing.T) {
	binRepo := infra.NewBinaryFileRepo(test.NewMockDb())
	binService := app.NewBinaryFileService(binRepo)
	form := editor.NewEntryForm(&MockEntry{}, binService)
	entry, err := form.Parse(NewMockFormData())
	require.NoError(t, err)
	require.NotZero(t, entry.MetaData().(*MockEntryMetaData).Image)
	require.Equal(t, "Content", entry.MetaData().(*MockEntryMetaData).Content)
}
