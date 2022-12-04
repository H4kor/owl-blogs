package web

import (
	"fmt"
	"h4kor/owl-blogs"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

func isUserLoggedIn(user *owl.User, r *http.Request) bool {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	return user.ValidateSession(sessionCookie.Value)
}

func setCSRFCookie(w http.ResponseWriter) string {
	csrfToken := owl.GenerateRandomString(32)
	cookie := http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	return csrfToken
}

func checkCSRF(r *http.Request) bool {
	// CSRF check
	formCsrfToken := r.FormValue("csrf_token")
	cookieCsrfToken, err := r.Cookie("csrf_token")

	if err != nil {
		println("Error getting csrf token from cookie: ", err.Error())
		return false
	}
	if formCsrfToken != cookieCsrfToken.Value {
		println("Invalid csrf token")
		return false
	}
	return true
}

func userLoginGetHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		if isUserLoggedIn(&user, r) {
			http.Redirect(w, r, user.EditorUrl(), http.StatusFound)
			return
		}
		csrfToken := setCSRFCookie(w)
		html, err := owl.RenderLoginPage(user, csrfToken)
		if err != nil {
			println("Error rendering login page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}
		w.Write([]byte(html))
	}
}

func userLoginPostHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}
		err = r.ParseForm()
		if err != nil {
			println("Error parsing form: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}

		// CSRF check
		if !checkCSRF(r) {
			w.WriteHeader(http.StatusBadRequest)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "CSRF Error",
				Message: "Invalid csrf token",
			})
			w.Write([]byte(html))
			return
		}

		password := r.Form.Get("password")
		if password == "" {
			userLoginGetHandler(repo)(w, r, ps)
			return
		}
		if !user.VerifyPassword(password) {
			userLoginGetHandler(repo)(w, r, ps)
			return
		}

		// set session cookie
		cookie := http.Cookie{
			Name:     "session",
			Value:    user.CreateNewSession(),
			Path:     "/",
			Expires:  time.Now().Add(30 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, user.EditorUrl(), http.StatusFound)
	}
}

func userEditorGetHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		if !isUserLoggedIn(&user, r) {
			http.Redirect(w, r, user.EditorLoginUrl(), http.StatusFound)
			return
		}

		csrfToken := setCSRFCookie(w)
		html, err := owl.RenderEditorPage(user, csrfToken)
		if err != nil {
			println("Error rendering editor page: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}
		w.Write([]byte(html))
	}
}

func userEditorPostHandler(repo *owl.Repository) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, err := getUserFromRepo(repo, ps)
		if err != nil {
			println("Error getting user: ", err.Error())
			notFoundHandler(repo)(w, r)
			return
		}

		if !isUserLoggedIn(&user, r) {
			http.Redirect(w, r, user.EditorLoginUrl(), http.StatusFound)
			return
		}

		err = r.ParseForm()
		if err != nil {
			println("Error parsing form: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}

		// CSRF check
		if !checkCSRF(r) {
			w.WriteHeader(http.StatusBadRequest)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "CSRF Error",
				Message: "Invalid csrf token",
			})
			w.Write([]byte(html))
			return
		}

		// get form values
		post_type := r.Form.Get("type")
		title := r.Form.Get("title")
		description := r.Form.Get("description")
		content := r.Form.Get("content")
		draft := r.Form.Get("draft")

		// conditional values
		reply_url := r.Form.Get("reply_url")
		bookmark_url := r.Form.Get("bookmark_url")

		// validate form values
		if post_type == "" {
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Missing post type",
				Message: "Post type is required",
			})
			w.Write([]byte(html))
			return
		}
		if (post_type == "article" || post_type == "page") && title == "" {
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Missing Title",
				Message: "Articles and Pages must have a title",
			})
			w.Write([]byte(html))
			return
		}
		if post_type == "reply" && reply_url == "" {
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Missing URL",
				Message: "You must provide a URL to reply to",
			})
			w.Write([]byte(html))
			return
		}
		if post_type == "bookmark" && bookmark_url == "" {
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Missing URL",
				Message: "You must provide a URL to bookmark",
			})
			w.Write([]byte(html))
			return
		}

		// TODO: scrape	reply_url for title and description
		// TODO: scrape bookmark_url for title and description

		// create post
		post, err := user.CreateNewPostFull(owl.PostMeta{
			Type:        post_type,
			Title:       title,
			Description: description,
			Draft:       draft == "on",
			Date:        time.Now(),
			Reply: owl.Reply{
				Url: reply_url,
			},
			Bookmark: owl.Bookmark{
				Url: bookmark_url,
			},
		}, content)

		if err != nil {
			println("Error creating post: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			html, _ := owl.RenderUserError(user, owl.ErrorMessage{
				Error:   "Internal server error",
				Message: "Internal server error",
			})
			w.Write([]byte(html))
			return
		}

		// redirect to post
		if !post.Meta().Draft {
			// scan for webmentions
			post.ScanForLinks()
			webmentions := post.OutgoingWebmentions()
			println("Found ", len(webmentions), " links")
			wg := sync.WaitGroup{}
			wg.Add(len(webmentions))
			for _, mention := range post.OutgoingWebmentions() {
				go func(mention owl.WebmentionOut) {
					fmt.Printf("Sending webmention to %s", mention.Target)
					defer wg.Done()
					post.SendWebmention(mention)
				}(mention)
			}
			wg.Wait()
			http.Redirect(w, r, post.FullUrl(), http.StatusFound)
		} else {
			http.Redirect(w, r, user.EditorUrl(), http.StatusFound)
		}
	}
}
