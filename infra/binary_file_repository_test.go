package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupBinaryRepo() repository.BinaryRepository {
	db := test.NewMockDb()
	repo := infra.NewBinaryFileRepo(db)
	return repo
}

func TestBinaryRepoCreate(t *testing.T) {
	repo := setupBinaryRepo()

	file, err := repo.Create("name", []byte("😀 😃 😄 😁"), nil)
	require.NoError(t, err)

	file, err = repo.FindById(file.Id)
	require.NoError(t, err)
	require.Equal(t, file.Name, "name")
	require.Equal(t, file.Data, []byte("😀 😃 😄 😁"))
}

func TestBinaryRepoNoSideEffect(t *testing.T) {
	repo := setupBinaryRepo()

	file, err := repo.Create("name1", []byte("111"), nil)
	require.NoError(t, err)

	file2, err := repo.Create("name2", []byte("222"), nil)
	require.NoError(t, err)

	file, err = repo.FindById(file.Id)
	require.NoError(t, err)
	file2, err = repo.FindById(file2.Id)
	require.NoError(t, err)
	require.Equal(t, file.Name, "name1")
	require.Equal(t, file.Data, []byte("111"))
	require.Equal(t, file2.Name, "name2")
	require.Equal(t, file2.Data, []byte("222"))
}

func TestBinaryWithSpaceInName(t *testing.T) {
	repo := setupBinaryRepo()

	file, err := repo.Create("name with space", []byte("111"), nil)
	require.NoError(t, err)

	file, err = repo.FindById(file.Id)
	require.NoError(t, err)
	require.Equal(t, file.Name, "name with space")
	require.Equal(t, file.Data, []byte("111"))
}
