package owl_test

import (
	"h4kor/owl-blogs"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type MockHtmlParser struct{}

func (*MockHtmlParser) ParseHEntry(resp *http.Response) (owl.ParsedHEntry, error) {
	return owl.ParsedHEntry{Title: "Mock Title"}, nil

}
func (*MockHtmlParser) ParseLinks(resp *http.Response) ([]string, error) {
	return []string{"http://example.com"}, nil

}
func (*MockHtmlParser) ParseLinksFromString(string) ([]string, error) {
	return []string{"http://example.com"}, nil

}
func (*MockHtmlParser) GetWebmentionEndpoint(resp *http.Response) (string, error) {
	return "http://example.com/webmention", nil

}

type MockHttpClient struct{}

func (*MockHttpClient) Get(url string) (resp *http.Response, err error) {
	return &http.Response{}, nil
}
func (*MockHttpClient) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {

	return &http.Response{}, nil
}
func (*MockHttpClient) PostForm(url string, data url.Values) (resp *http.Response, err error) {

	return &http.Response{}, nil
}

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

func randomUserName() string {
	return randomName()
}

func getTestUser() owl.User {
	repo, _ := owl.CreateRepository(testRepoName(), owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	return user
}

func getTestRepo(config owl.RepoConfig) owl.Repository {
	repo, _ := owl.CreateRepository(testRepoName(), config)
	return repo
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
