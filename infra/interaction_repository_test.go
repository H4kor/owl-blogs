package infra_test

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupInteractionRepo() repository.InteractionRepository {
	db := test.NewMockDb()
	register := app.NewInteractionTypeRegistry()
	register.Register(&test.MockInteraction{})
	repo := infra.NewInteractionRepo(db, register)
	return repo
}

func TestCreateInteraction(t *testing.T) {
	repo := setupInteractionRepo()
	i := &test.MockInteraction{}
	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i.SetEntryID("entryId")
	err := repo.Create(i)
	require.NoError(t, err)
	require.NotEmpty(t, i.ID())
}

func TestFindInteractionById(t *testing.T) {

	repo := setupInteractionRepo()
	i := &test.MockInteraction{}
	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i.SetEntryID("entryId")
	err := repo.Create(i)
	require.NoError(t, err)

	i2, err := repo.FindById(i.ID())
	require.NoError(t, err)
	require.Equal(t, i.ID(), i2.ID())
	require.Equal(t, i.Content(), i2.Content())
	meta := i.MetaData().(*test.MockInteractionMetaData)
	meta2 := i2.MetaData().(*test.MockInteractionMetaData)
	require.Equal(t, meta.Str, meta2.Str)
	require.Equal(t, meta.Number, meta2.Number)
	require.Equal(t, i2.EntryID(), "entryId")
}

func TestFindInteractionByEntryId(t *testing.T) {
	repo := setupInteractionRepo()
	i := &test.MockInteraction{}
	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i.SetEntryID("entryId")
	err := repo.Create(i)
	require.NoError(t, err)

	i2 := &test.MockInteraction{}
	i2.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i2.SetEntryID("entryId2")
	err = repo.Create(i2)
	require.NoError(t, err)

	inters, err := repo.FindAll("entryId")
	require.NoError(t, err)
	require.Equal(t, 1, len(inters))
}

func TestUpdateInteraction(t *testing.T) {
	repo := setupInteractionRepo()
	i := &test.MockInteraction{}
	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i.SetEntryID("entryId")
	err := repo.Create(i)
	require.NoError(t, err)

	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str2",
		Number: 2,
	})
	err = repo.Update(i)
	require.NoError(t, err)

	i2, err := repo.FindById(i.ID())
	require.NoError(t, err)
	meta := i.MetaData().(*test.MockInteractionMetaData)
	meta2 := i2.MetaData().(*test.MockInteractionMetaData)
	require.Equal(t, meta.Str, meta2.Str)
	require.Equal(t, meta.Number, meta2.Number)
}

func TestDeleteInteraction(t *testing.T) {
	repo := setupInteractionRepo()
	i := &test.MockInteraction{}
	i.SetMetaData(&test.MockInteractionMetaData{
		Str:    "str",
		Number: 1,
	})
	i.SetEntryID("entryId")
	err := repo.Create(i)
	require.NoError(t, err)

	err = repo.Delete(i)
	require.NoError(t, err)

	i2, err := repo.FindById(i.ID())
	require.Error(t, err)
	require.Nil(t, i2)
}
