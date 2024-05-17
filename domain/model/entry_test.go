package model_test

import (
	"owl-blogs/domain/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntryFullUrl(t *testing.T) {

	type testCase struct {
		Id   string
		Url  string
		Want string
	}

	testCases := []testCase{
		{Id: "foobar", Url: "https://example.com", Want: "https://example.com/posts/foobar/"},
		{Id: "foobar", Url: "https://example.com/", Want: "https://example.com/posts/foobar/"},
		{Id: "foobar", Url: "http://example.com", Want: "http://example.com/posts/foobar/"},
		{Id: "foobar", Url: "http://example.com/", Want: "http://example.com/posts/foobar/"},
		{Id: "bi-bar-buz", Url: "https://example.com", Want: "https://example.com/posts/bi-bar-buz/"},
		{Id: "foobar", Url: "https://example.com/lol/", Want: "https://example.com/lol/posts/foobar/"},
	}

	for _, test := range testCases {
		e := model.EntryBase{}
		e.SetID(test.Id)
		cfg := model.SiteConfig{FullUrl: test.Url}
		require.Equal(t, e.FullUrl(cfg), test.Want)
	}

}
