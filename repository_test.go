package owl_test

import (
	"h4kor/owl-blogs"
	"os"
	"path"
	"testing"
)

func TestCanCreateRepository(t *testing.T) {
	repoName := testRepoName()
	_, err := owl.CreateRepository(repoName)
	if err != nil {
		t.Error("Error creating repository: ", err.Error())
	}

}

func TestCannotCreateExistingRepository(t *testing.T) {
	repoName := testRepoName()
	owl.CreateRepository(repoName)
	_, err := owl.CreateRepository(repoName)
	if err == nil {
		t.Error("No error returned when creating existing repository")
	}
}

func TestCanCreateANewUser(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "")); err != nil {
		t.Error("User directory not created")
	}
}

func TestCannotRecreateExisitingUser(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	userName := randomUserName()
	repo.CreateUser(userName)
	_, err := repo.CreateUser(userName)
	if err == nil {
		t.Error("No error returned when creating existing user")
	}
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/VERSION")); err != nil {
		t.Error("Version file not created")
	}
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/base.html")); err != nil {
		t.Error("Base html file not created")
	}
}

func TestCreateUserAddConfigYml(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/config.yml")); err != nil {
		t.Error("Config file not created")
	}
}

func TestCreateUserAddsPublicFolder(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/public")); err != nil {
		t.Error("Public folder not created")
	}
}

func TestCanListRepoUsers(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user1, _ := repo.CreateUser(randomUserName())
	user2, _ := repo.CreateUser(randomUserName())
	// Create a new post
	users, _ := repo.Users()
	if len(users) != 2 {
		t.Error("Wrong number of users returned, expected 2, got ", len(users))
	}
	for _, user := range users {
		if user.Name() != user1.Name() && user.Name() != user2.Name() {
			t.Error("User found: " + user.Name())
		}
	}
}

func TestCanOpenRepository(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	repo, _ := owl.CreateRepository(repoName)
	// Open the repository
	repo2, err := owl.OpenRepository(repoName)
	if err != nil {
		t.Error("Error opening repository: ", err.Error())
	}
	if repo2.Dir() != repo.Dir() {
		t.Error("Repository directories do not match")
	}
}

func TestCannotOpenNonExisitingRepo(t *testing.T) {
	_, err := owl.OpenRepository(testRepoName())
	if err == nil {
		t.Error("No error returned when opening non-existing repository")
	}
}

func TestGetUser(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Get the user
	user2, err := repo.GetUser(user.Name())
	if err != nil {
		t.Error("Error getting user: ", err.Error())
	}
	if user2.Name() != user.Name() {
		t.Error("User names do not match")
	}
}

func TestCannotGetNonexistingUser(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	_, err := repo.GetUser(randomUserName())
	if err == nil {
		t.Error("No error returned when getting non-existing user")
	}
}

func TestGetStaticDirOfRepo(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	// Get the user
	staticDir := repo.StaticDir()
	if staticDir == "" {
		t.Error("Static dir not returned")
	}
}

func TestNewRepoGetsStaticFiles(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	if _, err := os.Stat(repo.StaticDir()); err != nil {
		t.Error("Static directory not found")
	}
	dir, _ := os.Open(repo.StaticDir())
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)

	if len(files) == 0 {
		t.Error("No static files found")
	}
}

func TestNewRepoGetsBaseHtml(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	if _, err := os.Stat(path.Join(repo.Dir(), "/base.html")); err != nil {
		t.Error("Base html file not found")
	}
}

func TestCanGetRepoTemplate(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	// Get the user
	template, err := repo.Template()
	if err != nil {
		t.Error("Error getting template: ", err.Error())
	}
	if template == "" {
		t.Error("Template not returned")
	}
}
