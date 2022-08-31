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
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "")); err != nil {
		t.Error("User directory not created")
	}
}

func TestCannotRecreateExisitingUser(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	userName := randomUserName()
	repo.CreateUser(userName)
	_, err := repo.CreateUser(userName)
	if err == nil {
		t.Error("No error returned when creating existing user")
	}
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/VERSION")); err != nil {
		t.Error("Version file not created")
	}
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/base.html")); err != nil {
		t.Error("Base html file not created")
	}
}

func TestCreateUserAddConfigYml(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/meta/config.yml")); err != nil {
		t.Error("Config file not created")
	}
}

func TestCreateUserAddsPublicFolder(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	if _, err := os.Stat(path.Join(user.Dir(), "/public")); err != nil {
		t.Error("Public folder not created")
	}
}

func TestCanListRepoUsers(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
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
	repo := getTestRepo()
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
	repo := getTestRepo()
	_, err := repo.GetUser(randomUserName())
	if err == nil {
		t.Error("No error returned when getting non-existing user")
	}
}

func TestGetStaticDirOfRepo(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	// Get the user
	staticDir := repo.StaticDir()
	if staticDir == "" {
		t.Error("Static dir not returned")
	}
}

func TestNewRepoGetsStaticFiles(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
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

func TestNewRepoGetsStaticFilesPicoCSSWithContent(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	file, err := os.Open(path.Join(repo.StaticDir(), "pico.min.css"))
	if err != nil {
		t.Error("Error opening pico.min.css")
	}
	// check that the file has content
	stat, _ := file.Stat()
	if stat.Size() == 0 {
		t.Error("pico.min.css is empty")
	}
}

func TestNewRepoGetsBaseHtml(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	if _, err := os.Stat(path.Join(repo.Dir(), "/base.html")); err != nil {
		t.Error("Base html file not found")
	}
}

func TestCanGetRepoTemplate(t *testing.T) {
	// Create a new user
	repo := getTestRepo()
	// Get the user
	template, err := repo.Template()
	if err != nil {
		t.Error("Error getting template: ", err.Error())
	}
	if template == "" {
		t.Error("Template not returned")
	}
}

func TestCanOpenRepositoryInSingleUserMode(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	userName := randomUserName()
	created_repo, _ := owl.CreateRepository(repoName)
	created_repo.CreateUser(userName)
	created_repo.CreateUser(randomUserName())
	created_repo.CreateUser(randomUserName())

	// Open the repository
	repo, _ := owl.OpenSingleUserRepo(repoName, userName)

	users, _ := repo.Users()
	if len(users) != 1 {
		t.Error("Wrong number of users returned, expected 1, got ", len(users))
	}
	if users[0].Name() != userName {
		t.Error("User name does not match")
	}
}

func TestSingleUserRepoUserUrlPathIsSimple(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	userName := randomUserName()
	created_repo, _ := owl.CreateRepository(repoName)
	created_repo.CreateUser(userName)

	// Open the repository
	repo, _ := owl.OpenSingleUserRepo(repoName, userName)
	user, _ := repo.GetUser(userName)
	if user.UrlPath() != "/" {
		t.Error("User url is not '/'. Got: ", user.UrlPath())
	}
}

func TestCanGetMapWithAllPostAliases(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	post, _ := user.CreateNewPost("test-1")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/bar\n"
	content += "  - /foo/baz\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	posts, _ := user.Posts()
	if len(posts) != 1 {
		t.Error("Wrong number of posts returned, expected 1, got ", len(posts))
	}

	var aliases map[string]*owl.Post
	aliases, err := repo.PostAliases()
	if err != nil {
		t.Error("Error getting post aliases: ", err.Error())
	}
	if len(aliases) != 2 {
		t.Error("Wrong number of aliases returned, expected 2, got ", len(aliases))
		t.Error("Aliases: ", aliases)
	}
	if aliases["/foo/bar"] == nil {
		t.Error("Alias '/foo/bar' not found")
	}
	if aliases["/foo/baz"] == nil {
		t.Error("Alias '/foo/baz' not found")
	}

}

func TestAliasesHaveCorrectPost(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser(randomUserName())
	post1, _ := user.CreateNewPost("test-1")
	post2, _ := user.CreateNewPost("test-2")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/1\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post1.ContentFile(), []byte(content), 0644)

	content = "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/2\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post2.ContentFile(), []byte(content), 0644)

	posts, _ := user.Posts()
	if len(posts) != 2 {
		t.Error("Wrong number of posts returned, expected 1, got ", len(posts))
	}

	var aliases map[string]*owl.Post
	aliases, err := repo.PostAliases()
	if err != nil {
		t.Error("Error getting post aliases: ", err.Error())
	}
	if len(aliases) != 2 {
		t.Error("Wrong number of aliases returned, expected 2, got ", len(aliases))
		t.Error("Aliases: ", aliases)
	}
	if aliases["/foo/1"].Id() != post1.Id() {
		t.Error("Alias '/foo/1' points to wrong post: ", aliases["/foo/1"].Id())
	}
	if aliases["/foo/2"].Id() != post2.Id() {
		t.Error("Alias '/foo/2' points to wrong post: ", aliases["/foo/2"].Id())
	}

}
