package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupAutherRepo() repository.AuthorRepository {
	db := test.NewMockDb()
	repo := infra.NewDefaultAuthorRepo(db)
	return repo
}

func TestAuthorRepoCreate(t *testing.T) {
	repo := setupAutherRepo()

	author, err := repo.Create("name", "password")
	require.NoError(t, err)

	author, err = repo.FindByName(author.Name)
	require.NoError(t, err)
	require.Equal(t, author.Name, "name")
	require.Equal(t, author.PasswordHash, "password")
}

func TestAuthorRepoNoSideEffect(t *testing.T) {
	repo := setupAutherRepo()

	author, err := repo.Create("name1", "password1")
	require.NoError(t, err)

	author2, err := repo.Create("name2", "password2")
	require.NoError(t, err)

	author, err = repo.FindByName(author.Name)
	require.NoError(t, err)
	author2, err = repo.FindByName(author2.Name)
	require.NoError(t, err)
	require.Equal(t, author.Name, "name1")
	require.Equal(t, author.PasswordHash, "password1")
	require.Equal(t, author2.Name, "name2")
	require.Equal(t, author2.PasswordHash, "password2")
}
