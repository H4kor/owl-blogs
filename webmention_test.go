package owl_test

import (
	"h4kor/owl-blogs"
	"testing"
)

//
// https://www.w3.org/TR/webmention/#h-webmention-verification
//

func TestParseValidHEntry(t *testing.T) {
	html := []byte("<div class=\"h-entry\"><div class=\"p-name\">Foo</div></div>")
	parser := &owl.OwlHtmlParser{}
	entry, err := parser.ParseHEntry(html)

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if entry.Title != "Foo" {
		t.Errorf("Wrong Title. Expected %v, got %v", "Foo", entry.Title)
	}
}

func TestParseValidHEntryWithoutTitle(t *testing.T) {
	html := []byte("<div class=\"h-entry\"></div><div class=\"p-name\">Foo</div>")
	parser := &owl.OwlHtmlParser{}
	entry, err := parser.ParseHEntry(html)

	if err != nil {
		t.Errorf("Unable to parse feed: %v", err)
	}
	if entry.Title != "" {
		t.Errorf("Wrong Title. Expected %v, got %v", "Foo", entry.Title)
	}
}
