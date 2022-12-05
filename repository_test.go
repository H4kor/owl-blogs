package owl_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/test/assertions"
	"os"
	"path"
	"testing"
)

func TestCanCreateRepository(t *testing.T) {
	repoName := testRepoName()
	_, err := owl.CreateRepository(repoName, owl.RepoConfig{})
	assertions.AssertNoError(t, err, "Error creating repository: ")

}

func TestCannotCreateExistingRepository(t *testing.T) {
	repoName := testRepoName()
	owl.CreateRepository(repoName, owl.RepoConfig{})
	_, err := owl.CreateRepository(repoName, owl.RepoConfig{})
	assertions.AssertError(t, err, "No error returned when creating existing repository")
}

func TestCanCreateANewUser(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	_, err := os.Stat(path.Join(user.Dir(), ""))
	assertions.AssertNoError(t, err, "Error creating user: ")
}

func TestCannotRecreateExisitingUser(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	userName := randomUserName()
	repo.CreateUser(userName)
	_, err := repo.CreateUser(userName)
	assertions.AssertError(t, err, "No error returned when creating existing user")
}

func TestCreateUserAddsVersionFile(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	_, err := os.Stat(path.Join(user.Dir(), "/meta/VERSION"))
	assertions.AssertNoError(t, err, "Version file not created")
}

func TestCreateUserAddsBaseHtmlFile(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	_, err := os.Stat(path.Join(user.Dir(), "/meta/base.html"))
	assertions.AssertNoError(t, err, "Base html file not created")
}

func TestCreateUserAddConfigYml(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	_, err := os.Stat(path.Join(user.Dir(), "/meta/config.yml"))
	assertions.AssertNoError(t, err, "Config file not created")
}

func TestCreateUserAddsPublicFolder(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	_, err := os.Stat(path.Join(user.Dir(), "/public"))
	assertions.AssertNoError(t, err, "Public folder not created")
}

func TestCanListRepoUsers(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user1, _ := repo.CreateUser(randomUserName())
	user2, _ := repo.CreateUser(randomUserName())
	// Create a new post
	users, _ := repo.Users()
	assertions.AssertLen(t, users, 2)
	for _, user := range users {
		assertions.AssertNot(
			t,
			user.Name() != user1.Name() && user.Name() != user2.Name(),
			"User found: "+user.Name(),
		)
	}
}

func TestCanOpenRepository(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	repo, _ := owl.CreateRepository(repoName, owl.RepoConfig{})
	// Open the repository
	repo2, err := owl.OpenRepository(repoName)
	assertions.AssertNoError(t, err, "Error opening repository: ")
	assertions.Assert(t, repo2.Dir() == repo.Dir(), "Repository directories do not match")
}

func TestCannotOpenNonExisitingRepo(t *testing.T) {
	_, err := owl.OpenRepository(testRepoName())
	assertions.AssertError(t, err, "No error returned when opening non-existing repository")
}

func TestGetUser(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Get the user
	user2, err := repo.GetUser(user.Name())
	assertions.AssertNoError(t, err, "Error getting user: ")
	assertions.Assert(t, user2.Name() == user.Name(), "User names do not match")
}

func TestCannotGetNonexistingUser(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	_, err := repo.GetUser(randomUserName())
	assertions.AssertError(t, err, "No error returned when getting non-existing user")
}

func TestGetStaticDirOfRepo(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	// Get the user
	staticDir := repo.StaticDir()
	assertions.Assert(t, staticDir != "", "Static dir is empty")
}

func TestNewRepoGetsStaticFiles(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	_, err := os.Stat(repo.StaticDir())
	assertions.AssertNoError(t, err, "Static dir not created")
	dir, _ := os.Open(repo.StaticDir())
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)

	assertions.AssertLen(t, files, 1)
}

func TestNewRepoGetsStaticFilesPicoCSSWithContent(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	file, err := os.Open(path.Join(repo.StaticDir(), "pico.min.css"))
	assertions.AssertNoError(t, err, "Error opening pico.min.css")
	// check that the file has content
	stat, _ := file.Stat()
	assertions.Assert(t, stat.Size() > 0, "pico.min.css is empty")
}

func TestNewRepoGetsBaseHtml(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	_, err := os.Stat(path.Join(repo.Dir(), "/base.html"))
	assertions.AssertNoError(t, err, "Base html file not found")
}

func TestCanGetRepoTemplate(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	// Get the user
	template, err := repo.Template()
	assertions.AssertNoError(t, err, "Error getting template: ")
	assertions.Assert(t, template != "", "Template is empty")
}

func TestCanOpenRepositoryInSingleUserMode(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	userName := randomUserName()
	created_repo, _ := owl.CreateRepository(repoName, owl.RepoConfig{SingleUser: userName})
	created_repo.CreateUser(userName)
	created_repo.CreateUser(randomUserName())
	created_repo.CreateUser(randomUserName())

	// Open the repository
	repo, _ := owl.OpenRepository(repoName)

	users, _ := repo.Users()
	assertions.AssertLen(t, users, 1)
	assertions.Assert(t, users[0].Name() == userName, "User name does not match")
}

func TestSingleUserRepoUserUrlPathIsSimple(t *testing.T) {
	// Create a new user
	repoName := testRepoName()
	userName := randomUserName()
	created_repo, _ := owl.CreateRepository(repoName, owl.RepoConfig{SingleUser: userName})
	created_repo.CreateUser(userName)

	// Open the repository
	repo, _ := owl.OpenRepository(repoName)
	user, _ := repo.GetUser(userName)
	assertions.Assert(t, user.UrlPath() == "/", "User url path is not /")
}

func TestCanGetMapWithAllPostAliases(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	post, _ := user.CreateNewPost(owl.PostMeta{Title: "test-1"}, "")

	content := "---\n"
	content += "title: Test\n"
	content += "aliases: \n"
	content += "  - /foo/bar\n"
	content += "  - /foo/baz\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	posts, _ := user.PublishedPosts()
	assertions.AssertLen(t, posts, 1)

	var aliases map[string]owl.IPost
	aliases, err := repo.PostAliases()
	assertions.AssertNoError(t, err, "Error getting post aliases: ")
	assertions.AssertMapLen(t, aliases, 2)
	assertions.Assert(t, aliases["/foo/bar"] != nil, "Alias '/foo/bar' not found")
	assertions.Assert(t, aliases["/foo/baz"] != nil, "Alias '/foo/baz' not found")

}

func TestAliasesHaveCorrectPost(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	post1, _ := user.CreateNewPost(owl.PostMeta{Title: "test-1"}, "")
	post2, _ := user.CreateNewPost(owl.PostMeta{Title: "test-2"}, "")

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

	posts, _ := user.PublishedPosts()
	assertions.AssertLen(t, posts, 2)

	var aliases map[string]owl.IPost
	aliases, err := repo.PostAliases()
	assertions.AssertNoError(t, err, "Error getting post aliases: ")
	assertions.AssertMapLen(t, aliases, 2)
	assertions.Assert(t, aliases["/foo/1"].Id() == post1.Id(), "Alias '/foo/1' does not point to post 1")
	assertions.Assert(t, aliases["/foo/2"].Id() == post2.Id(), "Alias '/foo/2' does not point to post 2")

}
