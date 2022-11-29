package owl_test

import (
	"fmt"
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/test/assertions"
	"os"
	"path"
	"testing"
)

func TestCreateNewPostCreatesEntryInPublic(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	files, err := os.ReadDir(path.Join(user.Dir(), "public"))
	assertions.AssertNoError(t, err, "Error reading directory")
	assertions.AssertLen(t, files, 1)
}

func TestCreateNewPostCreatesMediaDir(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost", false)
	_, err := os.Stat(post.MediaDir())
	assertions.AssertNot(t, os.IsNotExist(err), "Media directory not created")
}

func TestCreateNewPostAddsDateToMetaBlock(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost", false)
	posts, _ := user.PublishedPosts()
	post, _ := user.GetPost(posts[0].Id())
	meta := post.Meta()
	assertions.AssertNot(t, meta.Date.IsZero(), "Date not set")
}

func TestCreateNewPostMultipleCalls(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	files, err := os.ReadDir(path.Join(user.Dir(), "public"))
	assertions.AssertNoError(t, err, "Error reading directory")
	assertions.AssertEqual(t, len(files), 3)
}

func TestCanListUserPosts(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	posts, err := user.PublishedPosts()
	assertions.AssertNoError(t, err, "Error reading posts")
	assertions.AssertLen(t, posts, 3)
}

func TestCannotListUserPostsInSubdirectories(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	os.Mkdir(path.Join(user.PostDir(), "foo"), 0755)
	os.Mkdir(path.Join(user.PostDir(), "foo/bar"), 0755)
	content := ""
	content += "---\n"
	content += "title: test\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"

	os.WriteFile(path.Join(user.PostDir(), "foo/index.md"), []byte(content), 0644)
	os.WriteFile(path.Join(user.PostDir(), "foo/bar/index.md"), []byte(content), 0644)
	posts, _ := user.PublishedPosts()
	postIds := []string{}
	for _, p := range posts {
		postIds = append(postIds, p.Id())
	}
	if !contains(postIds, "foo") {
		t.Error("Does not contain post: foo. Found:")
		for _, p := range posts {
			t.Error("\t" + p.Id())
		}
	}

	if contains(postIds, "foo/bar") {
		t.Error("Invalid post found: foo/bar. Found:")
		for _, p := range posts {
			t.Error("\t" + p.Id())
		}
	}
}

func TestCannotListUserPostsWithoutIndexMd(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	os.Mkdir(path.Join(user.PostDir(), "foo"), 0755)
	os.Mkdir(path.Join(user.PostDir(), "foo/bar"), 0755)
	content := ""
	content += "---\n"
	content += "title: test\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"

	os.WriteFile(path.Join(user.PostDir(), "foo/bar/index.md"), []byte(content), 0644)
	posts, _ := user.PublishedPosts()
	postIds := []string{}
	for _, p := range posts {
		postIds = append(postIds, p.Id())
	}
	if contains(postIds, "foo") {
		t.Error("Contains invalid post: foo. Found:")
		for _, p := range posts {
			t.Error("\t" + p.Id())
		}
	}
}

func TestListUserPostsDoesNotIncludeDrafts(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost", false)
	content := ""
	content += "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	posts, _ := user.PublishedPosts()
	assertions.AssertLen(t, posts, 0)
}

func TestListUsersDraftsExcludedRealWorld(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost", false)
	content := ""
	content += "---\n"
	content += "title: Articles September 2019\n"
	content += "author: h4kor\n"
	content += "type: post\n"
	content += "date: -001-11-30T00:00:00+00:00\n"
	content += "draft: true\n"
	content += "url: /?p=426\n"
	content += "categories:\n"
	content += "  - Uncategorised\n"
	content += "\n"
	content += "---\n"
	content += "<https://nesslabs.com/time-anxiety>\n"

	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	posts, _ := user.PublishedPosts()
	assertions.AssertLen(t, posts, 0)
}

func TestCanLoadPost(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost", false)

	posts, _ := user.PublishedPosts()
	post, _ := user.GetPost(posts[0].Id())
	assertions.Assert(t, post.Title() == "testpost", "Post title is not correct")
}

func TestUserUrlPath(t *testing.T) {
	user := getTestUser()
	assertions.Assert(t, user.UrlPath() == "/user/"+user.Name()+"/", "Wrong url path")
}

func TestUserFullUrl(t *testing.T) {
	user := getTestUser()
	assertions.Assert(t, user.FullUrl() == "http://localhost:8080/user/"+user.Name()+"/", "Wrong url path")
}

func TestPostsSortedByPublishingDateLatestFirst(t *testing.T) {
	user := getTestUser()
	// Create a new post
	post1, _ := user.CreateNewPost("testpost", false)
	post2, _ := user.CreateNewPost("testpost2", false)

	content := "---\n"
	content += "title: Test Post\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post1.ContentFile(), []byte(content), 0644)

	content = "---\n"
	content += "title: Test Post 2\n"
	content += "date: Wed, 17 Aug 2022 20:50:06 +0000\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post2.ContentFile(), []byte(content), 0644)

	posts, _ := user.PublishedPosts()
	assertions.Assert(t, posts[0].Id() == post2.Id(), "Wrong Id")
	assertions.Assert(t, posts[1].Id() == post1.Id(), "Wrong Id")
}

func TestPostsSortedByPublishingDateLatestFirst2(t *testing.T) {
	user := getTestUser()
	// Create a new post
	posts := []*owl.Post{}
	for i := 59; i >= 0; i-- {
		post, _ := user.CreateNewPost("testpost", false)
		content := "---\n"
		content += "title: Test Post\n"
		content += fmt.Sprintf("date: Wed, 17 Aug 2022 10:%02d:02 +0000\n", i)
		content += "---\n"
		content += "This is a test"
		os.WriteFile(post.ContentFile(), []byte(content), 0644)
		posts = append(posts, post)
	}

	retPosts, _ := user.PublishedPosts()
	for i, p := range retPosts {
		assertions.Assert(t, p.Id() == posts[i].Id(), "Wrong Id")
	}
}

func TestPostsSortedByPublishingDateBrokenAtBottom(t *testing.T) {
	user := getTestUser()
	// Create a new post
	post1, _ := user.CreateNewPost("testpost", false)
	post2, _ := user.CreateNewPost("testpost2", false)

	content := "---\n"
	content += "title: Test Post\n"
	content += "date: Wed, 17 +0000\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post1.ContentFile(), []byte(content), 0644)

	content = "---\n"
	content += "title: Test Post 2\n"
	content += "date: Wed, 17 Aug 2022 20:50:06 +0000\n"
	content += "---\n"
	content += "This is a test"
	os.WriteFile(post2.ContentFile(), []byte(content), 0644)

	posts, _ := user.PublishedPosts()
	assertions.Assert(t, posts[0].Id() == post2.Id(), "Wrong Id")
	assertions.Assert(t, posts[1].Id() == post1.Id(), "Wrong Id")
}

func TestAvatarEmptyIfNotExist(t *testing.T) {
	user := getTestUser()
	assertions.Assert(t, user.AvatarUrl() == "", "Avatar should be empty")
}

func TestAvatarSetIfFileExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "avatar.png"), []byte("test"), 0644)
	assertions.Assert(t, user.AvatarUrl() != "", "Avatar should not be empty")
}

func TestPostNameIllegalFileName(t *testing.T) {
	user := getTestUser()
	_, err := user.CreateNewPost("testpost?///", false)
	assertions.AssertNoError(t, err, "Should not have failed")
}

func TestFaviconIfNotExist(t *testing.T) {
	user := getTestUser()
	assertions.Assert(t, user.FaviconUrl() == "", "Favicon should be empty")
}

func TestFaviconSetIfFileExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "favicon.ico"), []byte("test"), 0644)
	assertions.Assert(t, user.FaviconUrl() != "", "Favicon should not be empty")
}

func TestResetUserPassword(t *testing.T) {
	user := getTestUser()
	user.ResetPassword("test")
	assertions.Assert(t, user.Config().PassworHash != "", "Password Hash should not be empty")
	assertions.Assert(t, user.Config().PassworHash != "test", "Password Hash should not be test")
}

func TestVerifyPassword(t *testing.T) {
	user := getTestUser()
	user.ResetPassword("test")
	assertions.Assert(t, user.VerifyPassword("test"), "Password should be correct")
	assertions.Assert(t, !user.VerifyPassword("test2"), "Password should be incorrect")
	assertions.Assert(t, !user.VerifyPassword(""), "Password should be incorrect")
	assertions.Assert(t, !user.VerifyPassword("Test"), "Password should be incorrect")
	assertions.Assert(t, !user.VerifyPassword("TEST"), "Password should be incorrect")
	assertions.Assert(t, !user.VerifyPassword("0000000"), "Password should be incorrect")

}

func TestValidateAccessTokenWrongToken(t *testing.T) {
	user := getTestUser()
	code, _ := user.GenerateAuthCode(
		"test", "test", "test", "test", "test",
	)
	user.GenerateAccessToken(owl.AuthCode{
		Code:                code,
		ClientId:            "test",
		RedirectUri:         "test",
		CodeChallenge:       "test",
		CodeChallengeMethod: "test",
		Scope:               "test",
	})
	valid, _ := user.ValidateAccessToken("test")
	assertions.Assert(t, !valid, "Token should be invalid")
}

func TestValidateAccessTokenCorrectToken(t *testing.T) {
	user := getTestUser()
	code, _ := user.GenerateAuthCode(
		"test", "test", "test", "test", "test",
	)
	token, _, _ := user.GenerateAccessToken(owl.AuthCode{
		Code:                code,
		ClientId:            "test",
		RedirectUri:         "test",
		CodeChallenge:       "test",
		CodeChallengeMethod: "test",
		Scope:               "test",
	})
	valid, aToken := user.ValidateAccessToken(token)
	assertions.Assert(t, valid, "Token should be valid")
	assertions.Assert(t, aToken.ClientId == "test", "Token should be valid")
	assertions.Assert(t, aToken.Token == token, "Token should be valid")
}
