package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupSiteConfigRepo() repository.SiteConfigRepository {
	db := test.NewMockDb()
	repo := infra.NewSiteConfigRepo(db)
	return repo
}

func TestSiteConfigRepo(t *testing.T) {
	repo := setupSiteConfigRepo()

	config, err := repo.Get()
	require.NoError(t, err)
	require.Equal(t, "", config.Title)
	require.Equal(t, "", config.SubTitle)

	config.Title = "title"
	config.SubTitle = "SubTitle"

	err = repo.Update(config)
	require.NoError(t, err)

	config2, err := repo.Get()
	require.NoError(t, err)
	require.Equal(t, "title", config2.Title)
	require.Equal(t, "SubTitle", config2.SubTitle)
}

func TestSiteConfigUpdates(t *testing.T) {
	repo := setupSiteConfigRepo()

	config, err := repo.Get()
	require.NoError(t, err)
	require.Equal(t, "", config.Title)
	require.Equal(t, "", config.SubTitle)

	config.Title = "title"
	config.SubTitle = "SubTitle"

	err = repo.Update(config)
	require.NoError(t, err)

	config2, err := repo.Get()
	require.NoError(t, err)
	require.Equal(t, "title", config2.Title)
	require.Equal(t, "SubTitle", config2.SubTitle)

	config2.Title = "title2"
	config2.SubTitle = "SubTitle2"

	err = repo.Update(config2)
	require.NoError(t, err)

	config3, err := repo.Get()
	require.NoError(t, err)
	require.Equal(t, "title2", config3.Title)
	require.Equal(t, "SubTitle2", config3.SubTitle)

}
