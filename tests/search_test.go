package tests

import (
	"net/http/httptest"
	entrytypes "owl-blogs/entry_types"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/stretchr/testify/require"
)

func TestSearchHandling(t *testing.T) {
	type give struct {
		queryStr string
		notes []string
	}
	type want struct {
		contains string
		notContain string
	}
	tests := []struct{
		give
		want
	}{
		{
			give: give{queryStr: "", notes: []string{"a note"}},
			want: want{contains: "Search", notContain: "a note"},
		},
		{
			give: give{queryStr: "?query=note", notes: []string{"a note"}},
			want: want{contains: "a note", notContain: ""},
		},
		{
			give: give{queryStr: "?query=111", notes: []string{"111", "222"}},
			want: want{contains: "111", notContain: "222"},
		},
		{
			give: give{queryStr: "?query=Note", notes: []string{"a note"}},
			want: want{contains: "a note", notContain: ""},
		},
		{
			give: give{queryStr: "?query=AAA", notes: []string{"aaa", "aAa", "AAA"}},
			want: want{contains: "aAa", notContain: ""},
		},
		{
			give: give{queryStr: "?query=AAA", notes: []string{"aaa", "aAa", "AAA"}},
			want: want{contains: "aAa", notContain: ""},
		},
		{
			give: give{queryStr: "?query=AAA", notes: []string{"aaa", "aAa", "AAA"}},
			want: want{contains: "aaa", notContain: ""},
		},
		{
			give: give{queryStr: "?query=AAA", notes: []string{"aaa", "aAa", "AAA"}},
			want: want{contains: "AAA", notContain: ""},
		},
		{
			give: give{queryStr: "?query=h4", notes: []string{"#### a note"}},
			want: want{contains: "", notContain: "<h4>a note</h4>"},
		},
	}
	for _, test := range tests {
		// test
		app := DefaultTestApp()
		srv := adaptor.FiberApp(app.FiberApp)
		for _, nText := range test.give.notes {
			note := entrytypes.Note{}
			note.SetMetaData(&entrytypes.NoteMetaData{
			Content: nText,})
			now := time.Now()
			note.SetPublishedAt(&now)
			app.EntryService.Create(&note)
		}
		// test
		req := httptest.NewRequest("GET", "/search" + test.give.queryStr, nil)
		resp := httptest.NewRecorder()
		srv.ServeHTTP(resp, req)
		// validation
		if test.want.contains != "" {
			require.Contains(t, resp.Body.String(), test.want.contains)
		}
		if test.want.notContain != "" {
			require.NotContains(t, resp.Body.String(), test.want.notContain)
		}
	}
}

