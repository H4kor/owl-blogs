package owl_test

import (
	"fmt"
	"h4kor/owl-blogs"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestCreateNewPostCreatesEntryInPublic(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	files, err := ioutil.ReadDir(path.Join(user.Dir(), "public"))
	if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 1 {
		t.Error("Post not created")
	}
}

func TestCreateNewPostCreatesMediaDir(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost")
	if _, err := os.Stat(post.MediaDir()); os.IsNotExist(err) {
		t.Error("Media directory not created")
	}
}

func TestCreateNewPostAddsDateToMetaBlock(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost")
	posts, _ := user.Posts()
	post, _ := user.GetPost(posts[0].Id())
	meta := post.Meta()
	if meta.Date == "" {
		t.Error("Found no date. Got: " + meta.Date)
	}
}

func TestCreateNewPostMultipleCalls(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	files, err := ioutil.ReadDir(path.Join(user.Dir(), "public"))
	if err != nil {
		t.Error("Error reading directory")
	}
	if len(files) < 3 {
		t.Error(fmt.Sprintf("Only %d posts created", len(files)))
	}
}

func TestCanListUserPosts(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	user.CreateNewPost("testpost")
	posts, err := user.Posts()
	if err != nil {
		t.Error("Error reading posts")
	}
	if len(posts) != 3 {
		t.Error("No posts found")
	}
}

func TestCannotListUserPostsInSubdirectories(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
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
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	user.CreateNewPost("testpost")
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
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost")
	content := ""
	content += "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	posts, _ := user.Posts()
	if len(posts) != 0 {
		t.Error("Found draft post")
	}
}

func TestListUsersDraftsExcludedRealWorld(t *testing.T) {
	// Create a new user
	repo, _ := owl.CreateRepository(testRepoName())
	user, _ := repo.CreateUser(randomUserName())
	// Create a new post
	post, _ := user.CreateNewPost("testpost")
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
	if len(posts) != 0 {
		t.Error("Found draft post")
	}
}

func TestCanLoadPost(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost")

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
	post1, _ := user.CreateNewPost("testpost")
	post2, _ := user.CreateNewPost("testpost2")

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
		post, _ := user.CreateNewPost("testpost")
		content := "---\n"
		content += "title: Test Post\n"
		content += fmt.Sprintf("date: Wed, 17 Aug 2022 10:%02d:02 +0000\n", i)
		content += "---\n"
		content += "This is a test"
		os.WriteFile(post.ContentFile(), []byte(content), 0644)
		posts = append(posts, &post)
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
	post1, _ := user.CreateNewPost("testpost")
	post2, _ := user.CreateNewPost("testpost2")

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
