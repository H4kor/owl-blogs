package app_test

import (
	"owl-blogs/app"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupService() *app.EntryService {
	db := test.NewMockDb()
	register := app.NewEntryTypeRegistry()
	register.Register(&test.MockEntry{})
	repo := infra.NewEntryRepository(db, register)
	cfgRepo := infra.NewConfigRepo(db)
	cfgService := app.NewSiteConfigService(cfgRepo)
	service := app.NewEntryService(repo, cfgService, app.NewEventBus())
	return service
}

func TestNiceEntryId(t *testing.T) {
	service := setupService()
	entry := &test.MockEntry{}
	meta := test.MockEntryMetaData{
		Title: "Hello World",
	}
	entry.SetMetaData(&meta)

	err := service.Create(entry)
	require.NoError(t, err)
	require.Equal(t, "hello-world", entry.ID())
}

func TestPreserveSetEntryId(t *testing.T) {
	service := setupService()
	entry := &test.MockEntry{}
	meta := test.MockEntryMetaData{
		Title: "Hello World",
	}
	entry.SetMetaData(&meta)
	entry.SetID("foobar")

	err := service.Create(entry)
	require.NoError(t, err)
	require.Equal(t, "foobar", entry.ID())
}

func TestNoTitleCreation(t *testing.T) {
	service := setupService()
	entry := &test.MockEntry{}
	meta := test.MockEntryMetaData{
		Title: "",
	}
	entry.SetMetaData(&meta)

	err := service.Create(entry)
	require.NoError(t, err)
	require.NotEqual(t, "", entry.ID())
}
