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
	post, _ := user.GetPost(posts[0])
	_, meta := post.MarkdownData()
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
	if !contains(posts, "foo") {
		t.Error("Does not contain post: foo. Found:")
		for _, p := range posts {
			t.Error("\t" + p)
		}
	}

	if contains(posts, "foo/bar") {
		t.Error("Invalid post found: foo/bar. Found:")
		for _, p := range posts {
			t.Error("\t" + p)
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
	if contains(posts, "foo") {
		t.Error("Contains invalid post: foo. Found:")
		for _, p := range posts {
			t.Error("\t" + p)
		}
	}
}

func TestCanLoadPost(t *testing.T) {
	user := getTestUser()
	// Create a new post
	user.CreateNewPost("testpost")

	posts, _ := user.Posts()
	post, _ := user.GetPost(posts[0])
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
