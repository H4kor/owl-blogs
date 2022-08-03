package static_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/cmd/owl-web/static"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"
)

func randomName() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testRepoName() string {
	return "/tmp/" + randomName()
}

func getTestRepo() owl.Repository {
	repo, _ := owl.CreateRepository(testRepoName())
	return repo
}

func TestDeliversStaticFilesOfRepo(t *testing.T) {
	repo := getTestRepo()
	// create test static file
	fileName := "test.txt"
	filePath := path.Join(repo.StaticDir(), fileName)
	expected := "ok"
	err := os.WriteFile(filePath, []byte(expected), 0644)

	// Create Request and Response
	req, err := http.NewRequest("GET", "/static/test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.Handler(static.StaticHandler(repo))
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
