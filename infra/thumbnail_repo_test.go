package infra_test

import (
	"owl-blogs/app/repository"
	"owl-blogs/infra"
	"owl-blogs/test"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupThumbnailRepo() repository.ThumbnailRepository {
	db := test.NewMockDb()
	repo := infra.NewThumbnailRepo(db)
	return repo
}

func TestSave(t *testing.T) {
	repo := setupThumbnailRepo()

	thumb, err := repo.Save("foo.png", "image/png", []byte("test!!!"))
	require.NoError(t, err)

	require.Equal(t, thumb.BinaryFileId, "foo.png")
	require.Equal(t, thumb.MimeType, "image/png")
	require.Equal(t, thumb.Data, []byte("test!!!"))

}

func TestGet(t *testing.T) {
	repo := setupThumbnailRepo()

	thumb, err := repo.Save("foo.png", "image/png", []byte("test!!!"))
	require.NoError(t, err)

	thumb2, err := repo.Get("foo.png")
	require.NoError(t, err)

	require.Equal(t, thumb.BinaryFileId, thumb2.BinaryFileId)
	require.Equal(t, thumb.Data, thumb2.Data)
	require.Equal(t, thumb.MimeType, thumb2.MimeType)
}

func TestDelete(t *testing.T) {
	repo := setupThumbnailRepo()

	_, err := repo.Save("foo.png", "image/png", []byte("test!!!"))
	require.NoError(t, err)

	_, err = repo.Get("foo.png")
	require.NoError(t, err)

	err = repo.Delete("foo.png")
	require.NoError(t, err)

	_, err = repo.Get("foo.png")
	require.Error(t, err)
}
