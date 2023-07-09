package infra_test

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupRepo() repository.EntryRepository {
	db := test.NewMockDb()
	register := app.NewEntryTypeRegistry()
	register.Register(&test.MockEntry{})
	repo := infra.NewEntryRepository(db, register)
	return repo
}

func TestRepoCreate(t *testing.T) {
	repo := setupRepo()

	entry := &test.MockEntry{}
	now := time.Now()
	entry.SetPublishedAt(&now)
	entry.SetMetaData(&test.MockEntryMetaData{
		Str:    "str",
		Number: 1,
		Date:   now,
	})
	err := repo.Create(entry)
	require.NoError(t, err)

	entry2, err := repo.FindById(entry.ID())
	require.NoError(t, err)
	require.Equal(t, entry.ID(), entry2.ID())
	require.Equal(t, entry.Content(), entry2.Content())
	require.Equal(t, entry.PublishedAt().Unix(), entry2.PublishedAt().Unix())
	meta := entry.MetaData().(*test.MockEntryMetaData)
	meta2 := entry2.MetaData().(*test.MockEntryMetaData)
	require.Equal(t, meta.Str, meta2.Str)
	require.Equal(t, meta.Number, meta2.Number)
	require.Equal(t, meta.Date.Unix(), meta2.Date.Unix())
}

func TestRepoDelete(t *testing.T) {
	repo := setupRepo()

	entry := &test.MockEntry{}
	now := time.Now()
	entry.SetPublishedAt(&now)
	entry.SetMetaData(&test.MockEntryMetaData{
		Str:    "str",
		Number: 1,
		Date:   now,
	})
	err := repo.Create(entry)
	require.NoError(t, err)

	err = repo.Delete(entry)
	require.NoError(t, err)

	_, err = repo.FindById("id")
	require.Error(t, err)
}

func TestRepoFindAll(t *testing.T) {
	repo := setupRepo()

	entry := &test.MockEntry{}
	now := time.Now()
	entry.SetPublishedAt(&now)
	entry.SetMetaData(&test.MockEntryMetaData{
		Str:    "str",
		Number: 1,
		Date:   now,
	})
	err := repo.Create(entry)

	require.NoError(t, err)

	entry2 := &test.MockEntry{}
	now2 := time.Now()
	entry2.SetPublishedAt(&now2)
	entry2.SetMetaData(&test.MockEntryMetaData{
		Str:    "str",
		Number: 1,
		Date:   now,
	})

	err = repo.Create(entry2)
	require.NoError(t, err)

	entries, err := repo.FindAll(nil)
	require.NoError(t, err)
	require.Equal(t, 2, len(entries))

	entries, err = repo.FindAll(&[]string{"MockEntry"})
	require.NoError(t, err)
	require.Equal(t, 2, len(entries))

	entries, err = repo.FindAll(&[]string{"MockEntry2"})
	require.NoError(t, err)
	require.Equal(t, 0, len(entries))

}

func TestRepoUpdate(t *testing.T) {
	repo := setupRepo()

	entry := &test.MockEntry{}
	now := time.Now()
	entry.SetPublishedAt(&now)
	entry.SetMetaData(&test.MockEntryMetaData{
		Str:    "str",
		Number: 1,
		Date:   now,
	})
	err := repo.Create(entry)
	require.NoError(t, err)

	entry2 := &test.MockEntry{}
	now2 := time.Now()
	entry2.SetPublishedAt(&now2)
	entry2.SetMetaData(&test.MockEntryMetaData{
		Str:    "str2",
		Number: 2,
		Date:   now2,
	})
	err = repo.Create(entry2)
	require.NoError(t, err)
	err = repo.Update(entry2)
	require.NoError(t, err)

	entry3, err := repo.FindById(entry2.ID())
	require.NoError(t, err)
	require.Equal(t, entry3.Content(), entry2.Content())
	require.Equal(t, entry3.PublishedAt().Unix(), entry2.PublishedAt().Unix())
	meta := entry3.MetaData().(*test.MockEntryMetaData)
	meta2 := entry2.MetaData().(*test.MockEntryMetaData)
	require.Equal(t, meta.Str, meta2.Str)
	require.Equal(t, meta.Number, meta2.Number)
	require.Equal(t, meta.Date.Unix(), meta2.Date.Unix())
}

func TestRepoNoSideEffect(t *testing.T) {
	repo := setupRepo()

	entry1 := &test.MockEntry{}
	now1 := time.Now()
	entry1.SetPublishedAt(&now1)
	entry1.SetMetaData(&test.MockEntryMetaData{
		Str:    "1",
		Number: 1,
		Date:   now1,
	})

	err := repo.Create(entry1)
	require.NoError(t, err)

	entry2 := &test.MockEntry{}
	now2 := time.Now()
	entry2.SetPublishedAt(&now2)
	entry2.SetMetaData(&test.MockEntryMetaData{
		Str:    "2",
		Number: 2,
		Date:   now2,
	})
	err = repo.Create(entry2)
	require.NoError(t, err)

	r1, err := repo.FindById(entry1.ID())
	require.NoError(t, err)
	r2, err := repo.FindById(entry2.ID())
	require.NoError(t, err)

	require.Equal(t, r1.MetaData().(*test.MockEntryMetaData).Str, "1")
	require.Equal(t, r1.MetaData().(*test.MockEntryMetaData).Number, 1)
	require.Equal(t, r1.MetaData().(*test.MockEntryMetaData).Date.Unix(), now1.Unix())

	require.Equal(t, r2.MetaData().(*test.MockEntryMetaData).Str, "2")
	require.Equal(t, r2.MetaData().(*test.MockEntryMetaData).Number, 2)
	require.Equal(t, r2.MetaData().(*test.MockEntryMetaData).Date.Unix(), now2.Unix())

}
