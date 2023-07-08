package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testDbName() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 6)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return "/tmp/" + string(b) + ".db"
}

func TestEditorFormGet(t *testing.T) {
	db := infra.NewSqliteDB(testDbName())
	app := App(db).FiberApp

	req := httptest.NewRequest("GET", "/editor/ImageEntry", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
}

func TestEditorFormPost(t *testing.T) {
	dbName := testDbName()
	db := infra.NewSqliteDB(dbName)
	owlApp := App(db)
	app := owlApp.FiberApp
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

	req := httptest.NewRequest("POST", "/editor/ImageEntry", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 302, resp.StatusCode)
	require.Contains(t, resp.Header.Get("Location"), "/posts/")

	id := strings.Split(resp.Header.Get("Location"), "/")[2]
	entry, err := repo.FindById(id)
	require.NoError(t, err)
	require.Equal(t, "test content", entry.MetaData().(*model.ImageEntryMetaData).Content)
	imageId := entry.MetaData().(*model.ImageEntryMetaData).ImageId
	require.NotZero(t, imageId)
	bin, err := binRepo.FindById(imageId)
	require.NoError(t, err)
	require.Equal(t, bin.Name, "test.png")
	require.Equal(t, fileBytes, bin.Data)

}
