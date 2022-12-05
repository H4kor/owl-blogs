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

type Bookmark struct {
	Post
}

func (b *Bookmark) TemplateDir() string {
	return "bookmark"
}

type Reply struct {
	Post
}

func (r *Reply) TemplateDir() string {
	return "reply"
}
