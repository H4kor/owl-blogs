package owl

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

type User struct {
	repo *Repository
	name string
}

type UserConfig struct {
	Title       string   `yaml:"title"`
	SubTitle    string   `yaml:"subtitle"`
	HeaderColor string   `yaml:"header_color"`
	AuthorName  string   `yaml:"author_name"`
	Me          []UserMe `yaml:"me"`
	PassworHash string   `yaml:"password_hash"`
}

type UserMe struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type AuthCode struct {
	Code                string    `yaml:"code"`
	ClientId            string    `yaml:"client_id"`
	RedirectUri         string    `yaml:"redirect_uri"`
	CodeChallenge       string    `yaml:"code_challenge"`
	CodeChallengeMethod string    `yaml:"code_challenge_method"`
	Scope               string    `yaml:"scope"`
	Created             time.Time `yaml:"created"`
}

type AccessToken struct {
	Token       string    `yaml:"token"`
	Scope       string    `yaml:"scope"`
	ClientId    string    `yaml:"client_id"`
	RedirectUri string    `yaml:"redirect_uri"`
	Created     time.Time `yaml:"created"`
	ExpiresIn   int       `yaml:"expires_in"`
}

func (user User) Dir() string {
	return path.Join(user.repo.UsersDir(), user.name)
}

func (user User) UrlPath() string {
	return user.repo.UserUrlPath(user)
}

func (user User) FullUrl() string {
	url, _ := url.JoinPath(user.repo.FullUrl(), user.UrlPath())
	return url
}

func (user User) AuthUrl() string {
	if user.Config().PassworHash == "" {
		return ""
	}
	url, _ := url.JoinPath(user.FullUrl(), "auth/")
	return url
}

func (user User) TokenUrl() string {
	url, _ := url.JoinPath(user.AuthUrl(), "token/")
	return url
}

func (user User) IndieauthMetadataUrl() string {
	url, _ := url.JoinPath(user.FullUrl(), ".well-known/oauth-authorization-server")
	return url
}

func (user User) WebmentionUrl() string {
	url, _ := url.JoinPath(user.FullUrl(), "webmention/")
	return url
}

func (user User) MicropubUrl() string {
	url, _ := url.JoinPath(user.FullUrl(), "micropub/")
	return url
}

func (user User) MediaUrl() string {
	url, _ := url.JoinPath(user.UrlPath(), "media")
	return url
}

func (user User) PostDir() string {
	return path.Join(user.Dir(), "public")
}

func (user User) MetaDir() string {
	return path.Join(user.Dir(), "meta")
}

func (user User) MediaDir() string {
	return path.Join(user.Dir(), "media")
}

func (user User) ConfigFile() string {
	return path.Join(user.MetaDir(), "config.yml")
}

func (user User) AuthCodesFile() string {
	return path.Join(user.MetaDir(), "auth_codes.yml")
}

func (user User) AccessTokensFile() string {
	return path.Join(user.MetaDir(), "access_tokens.yml")
}

func (user User) Name() string {
	return user.name
}

func (user User) AvatarUrl() string {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif"} {
		if fileExists(path.Join(user.MediaDir(), "avatar"+ext)) {
			url, _ := url.JoinPath(user.MediaUrl(), "avatar"+ext)
			return url
		}
	}
	return ""
}

func (user User) FaviconUrl() string {
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".gif", ".ico"} {
		if fileExists(path.Join(user.MediaDir(), "favicon"+ext)) {
			url, _ := url.JoinPath(user.MediaUrl(), "favicon"+ext)
			return url
		}
	}
	return ""
}

func (user User) Posts() ([]*Post, error) {
	postFiles := listDir(path.Join(user.Dir(), "public"))
	posts := make([]*Post, 0)
	for _, id := range postFiles {
		// if is a directory and has index.md, add to posts
		if dirExists(path.Join(user.Dir(), "public", id)) {
			if fileExists(path.Join(user.Dir(), "public", id, "index.md")) {
				post, _ := user.GetPost(id)
				posts = append(posts, post)
			}
		}
	}

	// remove drafts
	n := 0
	for _, post := range posts {
		meta := post.Meta()
		if !meta.Draft {
			posts[n] = post
			n++
		}
	}
	posts = posts[:n]

	type PostWithDate struct {
		post *Post
		date time.Time
	}

	postDates := make([]PostWithDate, len(posts))
	for i, post := range posts {
		meta := post.Meta()
		postDates[i] = PostWithDate{post: post, date: meta.Date}
	}

	// sort posts by date
	sort.Slice(postDates, func(i, j int) bool {
		return postDates[i].date.After(postDates[j].date)
	})

	for i, post := range postDates {
		posts[i] = post.post
	}

	return posts, nil
}

func (user User) GetPost(id string) (*Post, error) {
	// check if posts index.md exists
	if !fileExists(path.Join(user.Dir(), "public", id, "index.md")) {
		return &Post{}, fmt.Errorf("post %s does not exist", id)
	}

	post := Post{user: &user, id: id}
	// post.loadMeta()
	meta := post.Meta()
	title := meta.Title
	post.title = fmt.Sprint(title)

	return &post, nil
}

func (user User) CreateNewPostFull(meta PostMeta, content string) (*Post, error) {
	slugHint := meta.Title
	if slugHint == "" {
		slugHint = "note"
	}
	folder_name := toDirectoryName(slugHint)
	post_dir := path.Join(user.Dir(), "public", folder_name)

	// if post already exists, add -n to the end of the name
	i := 0
	for {
		if dirExists(post_dir) {
			i++
			folder_name = toDirectoryName(fmt.Sprintf("%s-%d", slugHint, i))
			post_dir = path.Join(user.Dir(), "public", folder_name)
		} else {
			break
		}
	}
	post := Post{user: &user, id: folder_name, title: slugHint}

	initial_content := ""
	initial_content += "---\n"
	// write meta
	meta_bytes, err := yaml.Marshal(meta)
	if err != nil {
		return &Post{}, err
	}
	initial_content += string(meta_bytes)
	initial_content += "---\n"
	initial_content += "\n"
	initial_content += content

	// create post file
	os.Mkdir(post_dir, 0755)
	os.WriteFile(post.ContentFile(), []byte(initial_content), 0644)
	// create media dir
	os.Mkdir(post.MediaDir(), 0755)
	return &post, nil
}

func (user User) CreateNewPost(title string, draft bool) (*Post, error) {
	meta := PostMeta{
		Title:   title,
		Date:    time.Now(),
		Aliases: []string{},
		Draft:   draft,
	}
	return user.CreateNewPostFull(meta, title)
}

func (user User) Template() (string, error) {
	// load base.html
	path := path.Join(user.Dir(), "meta", "base.html")
	base_html, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(base_html), nil
}

func (user User) Config() UserConfig {
	meta := UserConfig{}
	loadFromYaml(user.ConfigFile(), &meta)
	return meta
}

func (user User) SetConfig(new_config UserConfig) error {
	return saveToYaml(user.ConfigFile(), new_config)
}

func (user User) PostAliases() (map[string]*Post, error) {
	post_aliases := make(map[string]*Post)
	posts, err := user.Posts()
	if err != nil {
		return post_aliases, err
	}
	for _, post := range posts {
		if err != nil {
			return post_aliases, err
		}
		for _, alias := range post.Aliases() {
			post_aliases[alias] = post
		}
	}
	return post_aliases, nil
}

func (user User) ResetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	config := user.Config()
	config.PassworHash = string(bytes)
	return user.SetConfig(config)
}

func (user User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Config().PassworHash), []byte(password),
	)
	return err == nil
}

func (user User) getAuthCodes() []AuthCode {
	codes := make([]AuthCode, 0)
	loadFromYaml(user.AuthCodesFile(), &codes)
	return codes
}

func (user User) addAuthCode(code AuthCode) error {
	codes := user.getAuthCodes()
	codes = append(codes, code)
	return saveToYaml(user.AuthCodesFile(), codes)
}

func (user User) GenerateAuthCode(
	client_id string, redirect_uri string,
	code_challenge string, code_challenge_method string,
	scope string,
) (string, error) {
	// generate code
	code := GenerateRandomString(32)
	return code, user.addAuthCode(AuthCode{
		Code:                code,
		ClientId:            client_id,
		RedirectUri:         redirect_uri,
		CodeChallenge:       code_challenge,
		CodeChallengeMethod: code_challenge_method,
		Scope:               scope,
		Created:             time.Now(),
	})
}

func (user User) VerifyAuthCode(
	code string, client_id string, redirect_uri string, code_verifier string,
) (bool, AuthCode) {
	codes := user.getAuthCodes()
	for _, c := range codes {
		if c.Code == code && c.ClientId == client_id && c.RedirectUri == redirect_uri {
			if c.CodeChallengeMethod == "plain" {
				return c.CodeChallenge == code_verifier, c
			} else if c.CodeChallengeMethod == "S256" {
				// hash code_verifier
				hash := sha256.Sum256([]byte(code_verifier))
				return c.CodeChallenge == base64.RawURLEncoding.EncodeToString(hash[:]), c
			} else if c.CodeChallengeMethod == "" {
				// Check age of code
				// A maximum lifetime of 10 minutes is recommended ( https://indieauth.spec.indieweb.org/#authorization-response)
				if time.Since(c.Created) < 10*time.Minute {
					return true, c
				}
			}
		}
	}
	return false, AuthCode{}
}

func (user User) getAccessTokens() []AccessToken {
	codes := make([]AccessToken, 0)
	loadFromYaml(user.AccessTokensFile(), &codes)
	return codes
}

func (user User) addAccessToken(code AccessToken) error {
	codes := user.getAccessTokens()
	codes = append(codes, code)
	return saveToYaml(user.AccessTokensFile(), codes)
}

func (user User) GenerateAccessToken(authCode AuthCode) (string, int, error) {
	// generate code
	token := GenerateRandomString(32)
	duration := 24 * 60 * 60
	return token, duration, user.addAccessToken(AccessToken{
		Token:       token,
		ClientId:    authCode.ClientId,
		RedirectUri: authCode.RedirectUri,
		Scope:       authCode.Scope,
		ExpiresIn:   duration,
		Created:     time.Now(),
	})
}

func (user User) ValidateAccessToken(token string) (bool, AccessToken) {
	tokens := user.getAccessTokens()
	for _, t := range tokens {
		if t.Token == token {
			return true, t
		}
	}
	return false, AccessToken{}
}
