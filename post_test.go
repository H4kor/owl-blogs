package owl_test

import (
	"h4kor/owl-blogs"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCanGetPostTitle(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result := post.Title()
	if result != "testpost" {
		t.Error("Wrong Title. Got: " + result)
	}
}

func TestMediaDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	result := post.MediaDir()
	if result != path.Join(post.Dir(), "media") {
		t.Error("Wrong MediaDir. Got: " + result)
	}
}

func TestPostUrlPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/"
	if !(post.UrlPath() == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestPostFullUrl(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	expected := "http://localhost:8080/user/" + user.Name() + "/posts/" + post.Id() + "/"
	if !(post.FullUrl() == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.FullUrl())
	}
}

func TestPostUrlMediaPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/data.png"
	if !(post.UrlMediaPath("data.png") == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestPostUrlMediaPathWithSubDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/foo/data.png"
	if !(post.UrlMediaPath("foo/data.png") == expected) {
		t.Error("Wrong url path")
		t.Error("Expected: " + expected)
		t.Error("     Got: " + post.UrlPath())
	}
}

func TestDraftInMetaData(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost("testpost", false)
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
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)
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
	repo := getTestRepo(owl.RepoConfig{AllowRawHtml: true})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)
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
	repo := getTestRepo(owl.RepoConfig{AllowRawHtml: true})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

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

	if post.Meta().Date.Format(time.RFC1123Z) != "Wed, 17 Aug 2022 10:50:02 +0000" {
		t.Errorf("Expected title: %s, got %s", "Wed, 17 Aug 2022 10:50:02 +0000", post.Meta().Title)
	}

	if post.Meta().Draft != true {
		t.Errorf("Expected title: %v, got %v", true, post.Meta().Draft)
	}
}

///
/// Webmention
///

func TestPersistIncomingWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)
	webmention := owl.WebmentionIn{
		Source: "http://example.com/source",
	}
	err := post.PersistIncomingWebmention(webmention)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	mentions := post.IncomingWebmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}

	if mentions[0].Source != webmention.Source {
		t.Errorf("Expected source: %s, got %s", webmention.Source, mentions[0].Source)
	}
}

func TestAddIncomingWebmentionCreatesFile(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	err := post.AddIncomingWebmention("https://example.com")
	if err != nil {
		t.Errorf("Got Error: %v", err)
	}

	mentions := post.IncomingWebmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}
}

func TestAddIncomingWebmentionNotOverwritingWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	post.PersistIncomingWebmention(owl.WebmentionIn{
		Source:         "https://example.com",
		ApprovalStatus: "approved",
	})

	post.AddIncomingWebmention("https://example.com")

	mentions := post.IncomingWebmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}

	if mentions[0].ApprovalStatus != "approved" {
		t.Errorf("Expected approval status: %s, got %s", "approved", mentions[0].ApprovalStatus)
	}
}

func TestEnrichAddsTitle(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	post.AddIncomingWebmention("https://example.com")
	post.EnrichWebmention(owl.WebmentionIn{Source: "https://example.com"})

	mentions := post.IncomingWebmentions()
	if len(mentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(mentions))
	}

	if mentions[0].Title != "Mock Title" {
		t.Errorf("Expected title: %s, got %s", "Mock Title", mentions[0].Title)
	}
}

func TestApprovedIncomingWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)
	webmention := owl.WebmentionIn{
		Source:         "http://example.com/source",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now(),
	}
	post.PersistIncomingWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source2",
		ApprovalStatus: "",
		RetrievedAt:    time.Now().Add(time.Hour * -1),
	}
	post.PersistIncomingWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source3",
		ApprovalStatus: "approved",
		RetrievedAt:    time.Now().Add(time.Hour * -2),
	}
	post.PersistIncomingWebmention(webmention)
	webmention = owl.WebmentionIn{
		Source:         "http://example.com/source4",
		ApprovalStatus: "rejected",
		RetrievedAt:    time.Now().Add(time.Hour * -3),
	}
	post.PersistIncomingWebmention(webmention)

	webmentions := post.ApprovedIncomingWebmentions()
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
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

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
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

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

func TestScanningForLinksDoesAddReplyUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "reply:\n"
	content += "  url: https://example.com/reply\n"
	content += "---\n"
	content += "\n"
	content += "Hi\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	post.ScanForLinks()
	webmentions := post.OutgoingWebmentions()
	if len(webmentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(webmentions))
	}
	if webmentions[0].Target != "https://example.com/reply" {
		t.Errorf("Expected target: %s, got %s", "https://example.com/reply", webmentions[0].Target)
	}
}

func TestCanSendWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

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

func TestSendWebmentionOnlyScansOncePerWeek(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	webmention := owl.WebmentionOut{
		Target:    "http://example.com",
		ScannedAt: time.Now().Add(time.Hour * -24 * 6),
	}

	post.PersistOutgoingWebmention(&webmention)
	webmentions := post.OutgoingWebmentions()
	webmention = webmentions[0]

	err := post.SendWebmention(webmention)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	webmentions = post.OutgoingWebmentions()

	if len(webmentions) != 1 {
		t.Errorf("Expected 1 webmention, got %d", len(webmentions))
	}

	if webmentions[0].ScannedAt != webmention.ScannedAt {
		t.Errorf("Expected ScannedAt to be unchanged. Expected: %v, got %v", webmention.ScannedAt, webmentions[0].ScannedAt)
	}
}

func TestSendingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 20; i++ {
		go func(k int) {
			webmention := owl.WebmentionOut{
				Target: "http://example.com" + strconv.Itoa(k),
			}
			post.SendWebmention(webmention)
			wg.Done()
		}(i)
	}

	wg.Wait()

	webmentions := post.OutgoingWebmentions()

	if len(webmentions) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(webmentions))
	}
}

func TestReceivingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	wg := sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 20; i++ {
		go func(k int) {
			post.AddIncomingWebmention("http://example.com" + strconv.Itoa(k))
			wg.Done()
		}(i)
	}

	wg.Wait()

	webmentions := post.IncomingWebmentions()

	if len(webmentions) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(webmentions))
	}

}

func TestSendingAndReceivingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	wg := sync.WaitGroup{}
	wg.Add(40)

	for i := 0; i < 20; i++ {
		go func(k int) {
			post.AddIncomingWebmention("http://example.com" + strconv.Itoa(k))
			wg.Done()
		}(i)
		go func(k int) {
			webmention := owl.WebmentionOut{
				Target: "http://example.com" + strconv.Itoa(k),
			}
			post.SendWebmention(webmention)
			wg.Done()
		}(i)
	}

	wg.Wait()

	ins := post.IncomingWebmentions()

	if len(ins) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(ins))
	}

	outs := post.OutgoingWebmentions()

	if len(outs) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(outs))
	}
}

func TestComplexParallelWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &MockHttpClient{}
	repo.Parser = &MockParseLinksHtmlParser{
		Links: []string{
			"http://example.com/1",
			"http://example.com/2",
			"http://example.com/3",
		},
	}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost("testpost", false)

	wg := sync.WaitGroup{}
	wg.Add(60)

	for i := 0; i < 20; i++ {
		go func(k int) {
			post.AddIncomingWebmention("http://example.com/" + strconv.Itoa(k))
			wg.Done()
		}(i)
		go func(k int) {
			webmention := owl.WebmentionOut{
				Target: "http://example.com/" + strconv.Itoa(k),
			}
			post.SendWebmention(webmention)
			wg.Done()
		}(i)
		go func() {
			post.ScanForLinks()
			wg.Done()
		}()
	}

	wg.Wait()

	ins := post.IncomingWebmentions()

	if len(ins) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(ins))
	}

	outs := post.OutgoingWebmentions()

	if len(outs) != 20 {
		t.Errorf("Expected 20 webmentions, got %d", len(outs))
	}
}

// func TestComplexParallelSimulatedProcessesWebmentions(t *testing.T) {
// 	repoName := testRepoName()
// 	repo, _ := owl.CreateRepository(repoName, owl.RepoConfig{})
// 	repo.HttpClient = &MockHttpClient{}
// 	repo.Parser = &MockParseLinksHtmlParser{
// 		Links: []string{
// 			"http://example.com/1",
// 			"http://example.com/2",
// 			"http://example.com/3",
// 		},
// 	}
// 	user, _ := repo.CreateUser("testuser")
// 	post, _ := user.CreateNewPost("testpost", false)

// 	wg := sync.WaitGroup{}
// 	wg.Add(40)

// 	for i := 0; i < 20; i++ {
// 		go func(k int) {
// 			defer wg.Done()
// 			fRepo, _ := owl.OpenRepository(repoName)
// 			fUser, _ := fRepo.GetUser("testuser")
// 			fPost, err := fUser.GetPost(post.Id())
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			fPost.AddIncomingWebmention("http://example.com/" + strconv.Itoa(k))
// 		}(i)
// 		go func(k int) {
// 			defer wg.Done()
// 			fRepo, _ := owl.OpenRepository(repoName)
// 			fUser, _ := fRepo.GetUser("testuser")
// 			fPost, err := fUser.GetPost(post.Id())
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			webmention := owl.WebmentionOut{
// 				Target: "http://example.com/" + strconv.Itoa(k),
// 			}
// 			fPost.SendWebmention(webmention)
// 		}(i)
// 	}

// 	wg.Wait()

// 	ins := post.IncomingWebmentions()

// 	if len(ins) != 20 {
// 		t.Errorf("Expected 20 webmentions, got %d", len(ins))
// 	}

// 	outs := post.OutgoingWebmentions()

// 	if len(outs) != 20 {
// 		t.Errorf("Expected 20 webmentions, got %d", len(outs))
// 	}
// }
