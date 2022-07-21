package kiss_test

import "h4kor/kiss-social"

func getTestUser() kiss.User {
	repo, _ := kiss.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	return user
}
