package owl_test

import (
	"os"
	"path"
	"testing"
)

func TestCanGetPostTitle(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result := post.Title()
	if result != "testpost" {
		t.Error("Wrong Title. Got: " + result)
	}
}

func TestMediaDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	result := post.MediaDir()
	if result != path.Join(post.Dir(), "media") {
		t.Error("Wrong MediaDir. Got: " + result)
	}
}

func TestPostUrlPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/"
	if !(post.UrlPath() == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestPostFullUrl(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	expected := "http://localhost:8080/user/" + user.Name() + "/posts/" + post.Id() + "/"
	if !(post.FullUrl() == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.FullUrl())
	}
}

func TestPostUrlMediaPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/data.png"
	if !(post.UrlMediaPath("data.png") == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestPostUrlMediaPathWithSubDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/foo/data.png"
	if !(post.UrlMediaPath("foo/data.png") == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestDraftInMetaData(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	_, meta := post.MarkdownData()
	if !meta.Draft {
		t.Error("Draft should be true")
	}

}
