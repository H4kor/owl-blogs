package owl

type Note struct {
	GenericPost
}

func (n *Note) TemplateDir() string {
	return "note"
}

type Article struct {
	GenericPost
}

func (a *Article) TemplateDir() string {
	return "article"
}

type Page struct {
	GenericPost
}

func (p *Page) TemplateDir() string {
	return "page"
}

type Bookmark struct {
	GenericPost
}

func (b *Bookmark) TemplateDir() string {
	return "bookmark"
}

type Reply struct {
	GenericPost
}

func (r *Reply) TemplateDir() string {
	return "reply"
}
