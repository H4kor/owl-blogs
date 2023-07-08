package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/test"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getUserToken(service *app.AuthorService) string {
	_, err := service.Create("test", "test")
	if err != nil {
		panic(err)
	}
	token, err := service.CreateToken("test")
	if err != nil {
		panic(err)
	}
	return token
}

func TestEditorFormGet(t *testing.T) {
	db := test.NewMockDb()
	owlApp := App(db)
	app := owlApp.FiberApp
	token := getUserToken(owlApp.AuthorService)

	req := httptest.NewRequest("GET", "/editor/Image", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestEditorFormGetNoAuth(t *testing.T) {
	db := test.NewMockDb()
	owlApp := App(db)
	app := owlApp.FiberApp

	req := httptest.NewRequest("GET", "/editor/Image", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid"})
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
}

func TestEditorFormPost(t *testing.T) {
	db := test.NewMockDb()
	owlApp := App(db)
	app := owlApp.FiberApp
	token := getUserToken(owlApp.AuthorService)
	repo := infra.NewEntryRepository(db, owlApp.Registry)
	binRepo := infra.NewBinaryFileRepo(db)

	fileDir, _ := os.Getwd()
	fileName := "../../test/fixtures/test.png"
	filePath := path.Join(fileDir, fileName)

	file, err := os.Open(filePath)
	require.NoError(t, err)
	fileBytes, err := ioutil.ReadFile(filePath)
	require.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("ImageId", filepath.Base(file.Name()))
	io.Copy(part, file)
	part, _ = writer.CreateFormField("Content")
	io.WriteString(part, "test content")
	writer.Close()

	req := httptest.NewRequest("POST", "/editor/Image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Location"), "/posts/")

	id := strings.Split(resp.Header.Get("Location"), "/")[2]
	entry, err := repo.FindById(id)
	require.NoError(t, err)
	require.Equal(t, "test content", entry.MetaData().(*model.ImageMetaData).Content)
	imageId := entry.MetaData().(*model.ImageMetaData).ImageId
	require.NotZero(t, imageId)
	bin, err := binRepo.FindById(imageId)
	require.NoError(t, err)
	require.Equal(t, bin.Name, "test.png")
	require.Equal(t, fileBytes, bin.Data)

}

func TestEditorFormPostNoAuth(t *testing.T) {
	db := test.NewMockDb()
	owlApp := App(db)
	app := owlApp.FiberApp

	fileDir, _ := os.Getwd()
	fileName := "../../test/fixtures/test.png"
	filePath := path.Join(fileDir, fileName)

	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("ImageId", filepath.Base(file.Name()))
	io.Copy(part, file)
	part, _ = writer.CreateFormField("Content")
	io.WriteString(part, "test content")
	writer.Close()

	req := httptest.NewRequest("POST", "/editor/Image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid"})
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Location"), "/auth/login")

}
