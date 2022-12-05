package owl_test

import (
	"h4kor/owl-blogs"
	"h4kor/owl-blogs/test/assertions"
	"h4kor/owl-blogs/test/mocks"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCanGetPostTitle(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	result := post.Title()
	assertions.AssertEqual(t, result, "testpost")
}

func TestMediaDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	result := post.MediaDir()
	assertions.AssertEqual(t, result, path.Join(post.Dir(), "media"))
}

func TestPostUrlPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/"
	assertions.AssertEqual(t, post.UrlPath(), expected)
}

func TestPostFullUrl(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	expected := "http://localhost:8080/user/" + user.Name() + "/posts/" + post.Id() + "/"
	assertions.AssertEqual(t, post.FullUrl(), expected)
}

func TestPostUrlMediaPath(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/data.png"
	assertions.AssertEqual(t, post.UrlMediaPath("data.png"), expected)
}

func TestPostUrlMediaPathWithSubDir(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	expected := "/user/" + user.Name() + "/posts/" + post.Id() + "/media/foo/data.png"
	assertions.AssertEqual(t, post.UrlMediaPath("foo/data.png"), expected)
}

func TestDraftInMetaData(t *testing.T) {
	user := getTestUser()
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "Write your post here.\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	meta := post.Meta()
	assertions.AssertEqual(t, meta.Draft, true)
}

func TestNoRawHTMLIfDisallowedByRepo(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "<script>alert('foo')</script>\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	html := post.RenderedContent()
	assertions.AssertNotContains(t, html, "<script>")
}

func TestRawHTMLIfAllowedByRepo(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{AllowRawHtml: true})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	content := "---\n"
	content += "title: test\n"
	content += "draft: true\n"
	content += "---\n"
	content += "\n"
	content += "<script>alert('foo')</script>\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)
	html := post.RenderedContent()
	assertions.AssertContains(t, html, "<script>")
}

func TestMeta(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{AllowRawHtml: true})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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

	assertions.AssertEqual(t, post.Meta().Title, "test")
	assertions.AssertLen(t, post.Meta().Aliases, 1)
	assertions.AssertEqual(t, post.Meta().Draft, true)
	assertions.AssertEqual(t, post.Meta().Date.Format(time.RFC1123Z), "Wed, 17 Aug 2022 10:50:02 +0000")
	assertions.AssertEqual(t, post.Meta().Draft, true)
}

///
/// Webmention
///

func TestPersistIncomingWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
	webmention := owl.WebmentionIn{
		Source: "http://example.com/source",
	}
	err := post.PersistIncomingWebmention(webmention)
	assertions.AssertNoError(t, err, "Error persisting webmention")
	mentions := post.IncomingWebmentions()
	assertions.AssertLen(t, mentions, 1)
	assertions.AssertEqual(t, mentions[0].Source, webmention.Source)
}

func TestAddIncomingWebmentionCreatesFile(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	err := post.AddIncomingWebmention("https://example.com")
	assertions.AssertNoError(t, err, "Error adding webmention")

	mentions := post.IncomingWebmentions()
	assertions.AssertLen(t, mentions, 1)
}

func TestAddIncomingWebmentionNotOverwritingWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	post.PersistIncomingWebmention(owl.WebmentionIn{
		Source:         "https://example.com",
		ApprovalStatus: "approved",
	})

	post.AddIncomingWebmention("https://example.com")

	mentions := post.IncomingWebmentions()
	assertions.AssertLen(t, mentions, 1)

	assertions.AssertEqual(t, mentions[0].ApprovalStatus, "approved")
}

func TestEnrichAddsTitle(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	post.AddIncomingWebmention("https://example.com")
	post.EnrichWebmention(owl.WebmentionIn{Source: "https://example.com"})

	mentions := post.IncomingWebmentions()
	assertions.AssertLen(t, mentions, 1)
	assertions.AssertEqual(t, mentions[0].Title, "Mock Title")
}

func TestApprovedIncomingWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")
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
	assertions.AssertLen(t, webmentions, 2)

	assertions.AssertEqual(t, webmentions[0].Source, "http://example.com/source")
	assertions.AssertEqual(t, webmentions[1].Source, "http://example.com/source3")

}

func TestScanningForLinks(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	content := "---\n"
	content += "title: test\n"
	content += "date: Wed, 17 Aug 2022 10:50:02 +0000\n"
	content += "---\n"
	content += "\n"
	content += "[Hello](https://example.com/hello)\n"
	os.WriteFile(post.ContentFile(), []byte(content), 0644)

	post.ScanForLinks()
	webmentions := post.OutgoingWebmentions()
	assertions.AssertLen(t, webmentions, 1)
	assertions.AssertEqual(t, webmentions[0].Target, "https://example.com/hello")
}

func TestScanningForLinksDoesNotAddDuplicates(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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
	assertions.AssertLen(t, webmentions, 1)
	assertions.AssertEqual(t, webmentions[0].Target, "https://example.com/hello")
}

func TestScanningForLinksDoesAddReplyUrl(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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
	assertions.AssertLen(t, webmentions, 1)
	assertions.AssertEqual(t, webmentions[0].Target, "https://example.com/reply")
}

func TestCanSendWebmention(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	webmention := owl.WebmentionOut{
		Target: "http://example.com",
	}

	err := post.SendWebmention(webmention)
	assertions.AssertNoError(t, err, "Error sending webmention")

	webmentions := post.OutgoingWebmentions()

	assertions.AssertLen(t, webmentions, 1)
	assertions.AssertEqual(t, webmentions[0].Target, "http://example.com")
	assertions.AssertEqual(t, webmentions[0].LastSentAt.IsZero(), false)
}

func TestSendWebmentionOnlyScansOncePerWeek(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

	webmention := owl.WebmentionOut{
		Target:    "http://example.com",
		ScannedAt: time.Now().Add(time.Hour * -24 * 6),
	}

	post.PersistOutgoingWebmention(&webmention)
	webmentions := post.OutgoingWebmentions()
	webmention = webmentions[0]

	err := post.SendWebmention(webmention)
	assertions.AssertError(t, err, "Expected error, got nil")

	webmentions = post.OutgoingWebmentions()

	assertions.AssertLen(t, webmentions, 1)
	assertions.AssertEqual(t, webmentions[0].ScannedAt, webmention.ScannedAt)
}

func TestSendingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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

	assertions.AssertLen(t, webmentions, 20)
}

func TestReceivingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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

	assertions.AssertLen(t, webmentions, 20)

}

func TestSendingAndReceivingMultipleWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockHtmlParser{}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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
	outs := post.OutgoingWebmentions()

	assertions.AssertLen(t, ins, 20)
	assertions.AssertLen(t, outs, 20)
}

func TestComplexParallelWebmentions(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	repo.HttpClient = &mocks.MockHttpClient{}
	repo.Parser = &mocks.MockParseLinksHtmlParser{
		Links: []string{
			"http://example.com/1",
			"http://example.com/2",
			"http://example.com/3",
		},
	}
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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
	outs := post.OutgoingWebmentions()

	assertions.AssertLen(t, ins, 20)
	assertions.AssertLen(t, outs, 20)
}

func TestPostWithoutContent(t *testing.T) {
	repo := getTestRepo(owl.RepoConfig{})
	user, _ := repo.CreateUser("testuser")
	post, _ := user.CreateNewPost(owl.PostMeta{}, "")

	result := post.RenderedContent()
	assertions.AssertEqual(t, result, "")
}

// func TestComplexParallelSimulatedProcessesWebmentions(t *testing.T) {
// 	repoName := testRepoName()
// 	repo, _ := owl.CreateRepository(repoName, owl.RepoConfig{})
// 	repo.HttpClient = &mocks.MockHttpClient{}
// 	repo.Parser = &MockParseLinksHtmlParser{
// 		Links: []string{
// 			"http://example.com/1",
// 			"http://example.com/2",
// 			"http://example.com/3",
// 		},
// 	}
// 	user, _ := repo.CreateUser("testuser")
// 	post, _ := user.CreateNewPostFull(owl.PostMeta{Type: "article", Title: "testpost"}, "")

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
