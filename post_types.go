package owl

type Note struct {
	Post
}

func (n *Note) TemplateDir() string {
	return "note"
}

type Article struct {
	Post
}

func (a *Article) TemplateDir() string {
	return "article"
}

type Page struct {
	Post
}

func (p *Page) TemplateDir() string {
	return "page"
}
