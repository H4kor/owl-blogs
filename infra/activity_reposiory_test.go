package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupActivityRepo() repository.ActivityRepository {
	db := test.NewMockDb()
	repo := infra.NewActivityRepo(db)
	return repo
}

func TestEmpty(t *testing.T) {
	repo := setupActivityRepo()
	acts, err := repo.ListRecent(0, 100)
	require.NoError(t, err)
	require.Len(t, acts, 0)
}

func TestInsertActivity(t *testing.T) {
	repo := setupActivityRepo()
	act := model.Activity{
		Id:        "https://example.com",
		Name:      "Test Name",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		Raw:       "foo",
		AuthorUrl: "https://example.com/user",
	}
	err := repo.Upsert(&act)
	require.NoError(t, err)

	acts, err := repo.ListRecent(0, 100)
	require.NoError(t, err)
	require.Len(t, acts, 1)
	require.Equal(t, act.Id, acts[0].Id)
	require.Equal(t, act.Name, acts[0].Name)
	require.Equal(t, act.Content, acts[0].Content)
	require.Equal(t, act.CreatedAt.Unix(), acts[0].CreatedAt.Unix())
	require.Equal(t, act.Raw, acts[0].Raw)
	require.Equal(t, act.AuthorUrl, acts[0].AuthorUrl)
}

func TestUpdateActivity(t *testing.T) {
	repo := setupActivityRepo()
	act := model.Activity{
		Id:        "https://example.com",
		Name:      "Test Name",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		Raw:       "foo",
		AuthorUrl: "https://example.com/user",
	}
	err := repo.Upsert(&act)
	require.NoError(t, err)
	err = repo.Upsert(&act)
	require.NoError(t, err)

	acts, err := repo.ListRecent(0, 100)
	require.NoError(t, err)
	require.Len(t, acts, 1)
	require.Equal(t, act.Id, acts[0].Id)
	require.Equal(t, act.Name, acts[0].Name)
	require.Equal(t, act.Content, acts[0].Content)
	require.Equal(t, act.CreatedAt.Unix(), acts[0].CreatedAt.Unix())
	require.Equal(t, act.Raw, acts[0].Raw)
	require.Equal(t, act.AuthorUrl, acts[0].AuthorUrl)
}
