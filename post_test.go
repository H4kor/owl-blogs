package owl_test

import (
	"h4kor/owl-blogs"
	"os"
	"path"
	"strings"
	"testing"
	"time"
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
	meta := post.Meta()
	if !meta.Draft {
		t.Error("Draft should be true")
	}

}

func TestNoRawHTMLIfDisallowedByRepo(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "<script>alert('foo')</script>\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	html := post.RenderedContent()
	html_str := html.String()
	if strings.Contains(html_str, "<script>") {
		t.Error("HTML should not be allowed")
	}
}

func TestRawHTMLIfAllowedByRepo(t *testing.T) {
	repo := getTestRepo()
	repo.SetAllowRawHtml(true)
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "<script>alert('foo')</script>\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	html := post.RenderedContent()
	html_str := html.String()
	if !strings.Contains(html_str, "<script>") {
		t.Error("HTML should be allowed")
	}
}

func TestLoadMeta(t *testing.T) {
	repo := getTestRepo()
	repo.SetAllowRawHtml(true)
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "aliases:\n"
	content += "  - foo/bar/\n"
	content += "---\n"
	content += "\n"
	content += "<script>alert('foo')</script>\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	err := post.LoadMeta()

	if err != nil {
		t.Errorf("Got Error: %v", err)
	}

	if post.Meta().Title != "test" {
		t.Errorf("Expected title: %s, got %s", "test", post.Meta().Title)
	}

	if len(post.Meta().Aliases) != 1 || post.Meta().Aliases[0] != "foo/bar/" {
		t.Errorf("Expected title: %v, got %v", []string{"foo/bar/"}, post.Meta().Aliases)
	}

	if post.Meta().Date != "Wed, 17 Aug 2022 10:50:02 +0000" {
		t.Errorf("Expected title: %s, got %s", "Wed, 17 Aug 2022 10:50:02 +0000", post.Meta().Title)
	}

	if post.Meta().Draft != true {
		t.Errorf("Expected title: %v, got %v", true, post.Meta().Draft)
	}
}

///
/// Webmention
///

func TestPersistWebmention(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")
	webmention := owl.WebmentionIn{
		Source: "http://example.com/source",
	}
	err := post.PersistWebmention(webmention)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	mentions := post.Webmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}

	if mentions[0].Source != webmention.Source {
		t.Errorf("Expected source: %s, got %s", webmention.Source, mentions[0].Source)
	}
}

func TestAddWebmentionCreatesFile(t *testing.T) {
	repo := getTestRepo()
	repo.HttpClient = &MockHttpRetriever{}
	repo.Parser = &MockHttpParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	err := post.AddWebmention("https://example.com")
	if err != nil {
		t.Errorf("Got Error: %v", err)
	}

	mentions := post.Webmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}
}

func TestAddWebmentionNotOverwritingFile(t *testing.T) {
	repo := getTestRepo()
	repo.HttpClient = &MockHttpRetriever{}
	repo.Parser = &MockHttpParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	post.AddWebmention("https://example.com")
	dir, _ := os.Open(post.WebmentionDir())
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)

	if len(files) != 1 {
		t.Error("No file created for webmention")
	}

	content := "url: https://example.com\n"
	content += "verified: true"
	os.WriteFile(path.Join(post.WebmentionDir(), files[0]), []byte(content), 0644)

	post.AddWebmention("https://example.com")

	fileContent, _ := os.ReadFile(path.Join(post.WebmentionDir(), files[0]))
	if string(fileContent) != content {
		t.Error("File content was modified.")
		t.Errorf("Got: %v", fileContent)
		t.Errorf("Expected: %v", content)
	}
}

func TestAddWebmentionAddsParsedTitle(t *testing.T) {
	repo := getTestRepo()
	repo.HttpClient = &MockHttpRetriever{}
	repo.Parser = &MockHttpParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	post.AddWebmention("https://example.com")
	dir, _ := os.Open(post.WebmentionDir())
	defer dir.Close()
	files, _ := dir.Readdirnames(-1)

	if len(files) != 1 {
		t.Error("No file created for webmention")
	}

	fileContent, _ := os.ReadFile(path.Join(post.WebmentionDir(), files[0]))
	if !strings.Contains(string(fileContent), "Mock Title") {
		t.Error("File not containing the title.")
		t.Errorf("Got: %v", string(fileContent))
	}
}

func TestApprovedWebmentions(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")
	webmention := owl.WebmentionIn{
		Source:         "http://example.com/source",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now(),
	}
	post.PersistWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source2",
		ApprovalStatus: "",
		RetrievedAt:    time.Now().Add(time.Hour * -1),
	}
	post.PersistWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source3",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now().Add(time.Hour * -2),
	}
	post.PersistWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source4",
		ApprovalStatus: "rejected",
		RetrievedAt:    time.Now().Add(time.Hour * -3),
	}
	post.PersistWebmention(webmention)

	webmentions := post.ApprovedWebmentions()
	if len(webmentions) != 2 {
		t.Errorf("Expected 2 webmentions, got %d", len(webmentions))
	}

	if webmentions[0].Source != "http://example.com/source" {
		t.Errorf("Expected source: %s, got %s", "http://example.com/source", webmentions[0].Source)
	}
	if webmentions[1].Source != "http://example.com/source3" {
		t.Errorf("Expected source: %s, got %s", "http://example.com/source3", webmentions[1].Source)
	}

}

func TestScanningForLinks(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "---\n"
	content += "\n"
	content += "[Hello](https://example.com/hello)\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	post.ScanForLinks()
	webmentions := post.OutgoingWebmentions()
	if len(webmentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(webmentions))
	}
	if webmentions[0].Target != "https://example.com/hello" {
		t.Errorf("Expected target: %s, got %s", "https://example.com/hello", webmentions[0].Target)
	}
}

func TestScanningForLinksDoesNotAddDuplicates(t *testing.T) {
	repo := getTestRepo()
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "---\n"
	content += "\n"
	content += "[Hello](https://example.com/hello)\n"
	content += "[Hello](https://example.com/hello)\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	post.ScanForLinks()
	post.ScanForLinks()
	post.ScanForLinks()
	webmentions := post.OutgoingWebmentions()
	if len(webmentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(webmentions))
	}
	if webmentions[0].Target != "https://example.com/hello" {
		t.Errorf("Expected target: %s, got %s", "https://example.com/hello", webmentions[0].Target)
	}
}

func TestCanSendWebmention(t *testing.T) {
	repo := getTestRepo()
	repo.HttpClient = &MockHttpRetriever{}
	repo.Parser = &MockHttpParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost")

	webmention := owl.WebmentionOut{
		Target: "http://example.com",
	}

	err := post.SendWebmention(webmention)
	if err != nil {
		t.Errorf("Error sending webmention: %v", err)
	}

	webmentions := post.OutgoingWebmentions()

	if len(webmentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(webmentions))
	}

	if webmentions[0].Target != "http://example.com" {
		t.Errorf("Expected target: %s, got %s", "http://example.com", webmentions[0].Target)
	}

	if webmentions[0].LastSentAt.IsZero() {
		t.Errorf("Expected LastSentAt to be set")
	}
}
