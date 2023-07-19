package app_test

import (
	"owl-blogs/app"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/test"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConfigRepo struct {
	config model.SiteConfig
}

// Get implements repository.SiteConfigRepository.
func (c *testConfigRepo) Get(name string, result interface{}) error {
	*result.(*model.SiteConfig) = c.config
	return nil
}

// Update implements repository.SiteConfigRepository.
func (c *testConfigRepo) Update(name string, result interface{}) error {
	c.config = result.(model.SiteConfig)
	return nil
}

func getAutherService() *app.AuthorService {
	db := test.NewMockDb()
	authorRepo := infra.NewDefaultAuthorRepo(db)
	authorService := app.NewAuthorService(authorRepo, &testConfigRepo{})
	return authorService

}

func TestAuthorCreate(t *testing.T) {
	authorService := getAutherService()
	author, err := authorService.Create("test", "test")
	require.NoError(t, err)
	require.Equal(t, "test", author.Name)
	require.NotEmpty(t, author.PasswordHash)
	require.NotEqual(t, "test", author.PasswordHash)
}

func TestAuthorFindByName(t *testing.T) {
	authorService := getAutherService()
	_, err := authorService.Create("test", "test")
	require.NoError(t, err)
	author, err := authorService.FindByName("test")
	require.NoError(t, err)
	require.Equal(t, "test", author.Name)
	require.NotEmpty(t, author.PasswordHash)
	require.NotEqual(t, "test", author.PasswordHash)
}

func TestAuthorAuthenticate(t *testing.T) {
	authorService := getAutherService()
	_, err := authorService.Create("test", "test")
	require.NoError(t, err)
	require.True(t, authorService.Authenticate("test", "test"))
	require.False(t, authorService.Authenticate("test", "test1"))
	require.False(t, authorService.Authenticate("test1", "test"))
}

func TestAuthorCreateToken(t *testing.T) {
	authorService := getAutherService()
	_, err := authorService.Create("test", "test")
	require.NoError(t, err)
	token, err := authorService.CreateToken("test")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEqual(t, "test", token)
}

func TestAuthorValidateToken(t *testing.T) {
	authorService := getAutherService()
	_, err := authorService.Create("test", "test")
	require.NoError(t, err)
	token, err := authorService.CreateToken("test")
	require.NoError(t, err)

	valid, name := authorService.ValidateToken(token)
	require.True(t, valid)
	require.Equal(t, "test", name)
	valid, _ = authorService.ValidateToken(token[:len(token)-2])
	require.False(t, valid)
	valid, _ = authorService.ValidateToken("test")
	require.False(t, valid)
	valid, _ = authorService.ValidateToken("test.test")
	require.False(t, valid)
	valid, _ = authorService.ValidateToken(strings.Replace(token, "test", "test1", 1))
	require.False(t, valid)
}
