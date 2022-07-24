package kiss_test

import (
	"h4kor/kiss-social"
	"math/rand"
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

func randomUserName() string {
	return randomName()
}

func getTestUser() kiss.User {
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	return user
}

func getTestRepo() kiss.Repository {
	repo, _ := kiss.CreateRepository(testRepoName())
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
