package owl_test

import (
	"fmt"
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/priv/assertions"
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
	posts, _ := user.Posts()
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
	if len(files) < 3 {
		t.Errorf("Only %d posts created", len(files))
	}
}

func TestCanListUserPosts(t *testing.T) {
	// Create a new user
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	user.CreateNewPost("testpost", false)
	posts, err := user.Posts()
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
	posts, _ := user.Posts()
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
	posts, _ := user.Posts()
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

	posts, _ := user.Posts()
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

	posts, _ := user.Posts()
	assertions.AssertLen(t, posts, 0)
}

func TestCanLoadPost(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost", false)

	posts, _ := user.Posts()
	post, _ := user.GetPost(posts[0].Id())
	if post.Title() != "testpost" {
		t.Error("Wrong title, Got: " + post.Title())
	}
}

func TestUserUrlPath(t *testing.T) {
	user := getTestUser()
	if !(user.UrlPath() == "/user/"+user.Name()+"/") {
		t.Error("Wrong url path, Expected: " + "/user/" + user.Name() + "/" + " Got: " + user.UrlPath())
	}
}

func TestUserFullUrl(t *testing.T) {
	user := getTestUser()
	if !(user.FullUrl() == "http://localhost:8080/user/"+user.Name()+"/") {
		t.Error("Wrong url path, Expected: " + "http://localhost:8080/user/" + user.Name() + "/" + " Got: " + user.FullUrl())
	}
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

	posts, _ := user.Posts()
	if posts[0].Id() != post2.Id() {
		t.Error("Wrong Id, Got: " + posts[0].Id())
	}
	if posts[1].Id() != post1.Id() {
		t.Error("Wrong Id, Got: " + posts[1].Id())
	}
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

	retPosts, _ := user.Posts()
	for i, p := range retPosts {
		if p.Id() != posts[i].Id() {
			t.Error("Wrong Id, Got: " + p.Id())
		}
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

	posts, _ := user.Posts()
	if posts[0].Id() != post2.Id() {
		t.Error("Wrong Id, Got: " + posts[0].Id())
	}
	if posts[1].Id() != post1.Id() {
		t.Error("Wrong Id, Got: " + posts[1].Id())
	}
}

func TestAvatarEmptyIfNotExist(t *testing.T) {
	user := getTestUser()
	if user.AvatarUrl() != "" {
		t.Error("Avatar should be empty")
	}
}

func TestAvatarSetIfFileExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "avatar.png"), []byte("test"), 0644)
	if user.AvatarUrl() == "" {
		t.Error("Avatar should not be empty")
	}
}

func TestPostNameIllegalFileName(t *testing.T) {
	user := getTestUser()
	_, err := user.CreateNewPost("testpost?///", false)
	assertions.AssertNoError(t, err, "Should not have failed")
}

func TestFaviconIfNotExist(t *testing.T) {
	user := getTestUser()
	if user.FaviconUrl() != "" {
		t.Error("Favicon should be empty")
	}
}

func TestFaviconSetIfFileExist(t *testing.T) {
	user := getTestUser()
	os.WriteFile(path.Join(user.MediaDir(), "favicon.ico"), []byte("test"), 0644)
	if user.FaviconUrl() == "" {
		t.Error("Favicon should not be empty")
	}
}
