package model

import "html/template"

type MeLinks struct {
	Name string
	Url  string
}

type EntryList struct {
	Id       string
	Title    string
	Include  []string
	ListType string
}

type MenuItem struct {
	Title string
	List  string
	Url   string
	Post  string
}

type SiteConfig struct {
	Title              string
	SubTitle           string
	PrimaryColor       string
	AuthorName         string
	Me                 []MeLinks
	Lists              []EntryList
	PrimaryListInclude []string
	HeaderMenu         []MenuItem
	FooterMenu         []MenuItem
	Secret             string
	AvatarUrl          string
	FullUrl            string
	HtmlHeadExtra      template.HTML
	FooterExtra        template.HTML
}
