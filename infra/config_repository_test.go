package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupSiteConfigRepo() repository.ConfigRepository {
	db := test.NewMockDb()
	repo := infra.NewConfigRepo(db)
	return repo
}

func TestSiteConfigRepo(t *testing.T) {
	repo := setupSiteConfigRepo()

	config := model.SiteConfig{}
	err := repo.Get("test", &config)
	require.NoError(t, err)
	require.Equal(t, "", config.Title)
	require.Equal(t, "", config.SubTitle)

	config.Title = "title"
	config.SubTitle = "SubTitle"

	err = repo.Update("test", config)
	require.NoError(t, err)

	config2 := model.SiteConfig{}
	err = repo.Get("test", &config2)
	require.NoError(t, err)
	require.Equal(t, "title", config2.Title)
	require.Equal(t, "SubTitle", config2.SubTitle)
}

func TestSiteConfigUpdates(t *testing.T) {
	repo := setupSiteConfigRepo()
	config := model.SiteConfig{}
	err := repo.Get("test", &config)
	require.NoError(t, err)
	require.Equal(t, "", config.Title)
	require.Equal(t, "", config.SubTitle)

	config.Title = "title"
	config.SubTitle = "SubTitle"

	err = repo.Update("test", config)
	require.NoError(t, err)
	config2 := model.SiteConfig{}
	err = repo.Get("test", &config2)
	require.NoError(t, err)
	require.Equal(t, "title", config2.Title)
	require.Equal(t, "SubTitle", config2.SubTitle)

	config2.Title = "title2"
	config2.SubTitle = "SubTitle2"

	err = repo.Update("test", config2)
	require.NoError(t, err)
	config3 := model.SiteConfig{}
	err = repo.Get("test", &config3)
	require.NoError(t, err)
	require.Equal(t, "title2", config3.Title)
	require.Equal(t, "SubTitle2", config3.SubTitle)

}
