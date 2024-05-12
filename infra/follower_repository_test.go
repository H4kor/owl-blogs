package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupFollowerRepo() repository.FollowerRepository {
	db := test.NewMockDb()
	repo := infra.NewFollowerRepository(db)
	return repo
}

func TestAddFollower(t *testing.T) {
	repo := setupFollowerRepo()

	err := repo.Add("foo@example.com")
	require.NoError(t, err)

	followers, err := repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 1)
	require.Equal(t, followers[0], "foo@example.com")
}

func TestDoubleAddFollower(t *testing.T) {
	repo := setupFollowerRepo()

	err := repo.Add("foo@example.com")
	require.NoError(t, err)

	err = repo.Add("foo@example.com")
	require.NoError(t, err)

	followers, err := repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 1)
	require.Equal(t, followers[0], "foo@example.com")
}

func TestMultipleAddFollower(t *testing.T) {
	repo := setupFollowerRepo()

	err := repo.Add("foo@example.com")
	require.NoError(t, err)

	err = repo.Add("bar@example.com")
	require.NoError(t, err)

	err = repo.Add("baz@example.com")
	require.NoError(t, err)

	followers, err := repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 3)
}

func TestRemoveFollower(t *testing.T) {
	repo := setupFollowerRepo()

	err := repo.Add("foo@example.com")
	require.NoError(t, err)

	followers, err := repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 1)

	err = repo.Remove("foo@example.com")
	require.NoError(t, err)

	followers, err = repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 0)

}

func TestRemoveNonExistingFollower(t *testing.T) {
	repo := setupFollowerRepo()

	err := repo.Remove("foo@example.com")
	require.NoError(t, err)

	followers, err := repo.All()
	require.NoError(t, err)
	require.Len(t, followers, 0)

}
