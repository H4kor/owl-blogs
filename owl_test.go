package owl_test

import (
	"h4kor/owl-blogs"
	"math/rand"
	"net/url"
	"time"
)

type MockHttpParser struct{}

func (*MockHttpParser) ParseHEntry(data []byte) (owl.ParsedHEntry, error) {
	return owl.ParsedHEntry{Title: "Mock Title"}, nil
}

func (*MockHttpParser) ParseLinks(data []byte) ([]string, error) {
	return []string{"http://example.com"}, nil
}

func (*MockHttpParser) GetWebmentionEndpoint(data []byte) (string, error) {
	return "http://example.com/webmention", nil
}

type MockHttpRetriever struct{}

func (*MockHttpRetriever) Get(url string) ([]byte, error) {
	return []byte(""), nil
}

func (m *MockHttpRetriever) Post(url string, data url.Values) ([]byte, error) {
	return []byte(""), nil
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
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	return user
}

func getTestRepo() owl.Repository {
	repo, _ := owl.CreateRepository(testRepoName())
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
