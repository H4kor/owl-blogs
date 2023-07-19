package importer

type V1UserConfig struct {
	Title              string       `yaml:"title"`
	SubTitle           string       `yaml:"subtitle"`
	HeaderColor        string       `yaml:"header_color"`
	AuthorName         string       `yaml:"author_name"`
	Me                 []V1UserMe   `yaml:"me"`
	PassworHash        string       `yaml:"password_hash"`
	Lists              []V1PostList `yaml:"lists"`
	PrimaryListInclude []string     `yaml:"primary_list_include"`
	HeaderMenu         []V1MenuItem `yaml:"header_menu"`
	FooterMenu         []V1MenuItem `yaml:"footer_menu"`
}

type V1UserMe struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type V1PostList struct {
	Id       string   `yaml:"id"`
	Title    string   `yaml:"title"`
	Include  []string `yaml:"include"`
	ListType string   `yaml:"list_type"`
}

type V1MenuItem struct {
	Title string `yaml:"title"`
	List  string `yaml:"list"`
	Url   string `yaml:"url"`
	Post  string `yaml:"post"`
}
