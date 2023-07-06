package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEditorFormGet(t *testing.T) {
	app := App().FiberApp

	req := httptest.NewRequest("GET", "/editor/ImageEntry", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestEditorFormPost(t *testing.T) {
	app := App().FiberApp

	fileDir, _ := os.Getwd()
	fileName := "../../test/fixtures/test.png"
	filePath := path.Join(fileDir, fileName)

	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("ImagePath", filepath.Base(file.Name()))
	io.Copy(part, file)
	part, _ = writer.CreateFormField("Content")
	io.WriteString(part, "test content")
	writer.Close()

	req := httptest.NewRequest("POST", "/editor/ImageEntry", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}
