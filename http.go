package owl

import (
	"io"
	"net/http"
	"net/url"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

type OwlHttpClient = http.Client
