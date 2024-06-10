package infra

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"

	"github.com/jmoiron/sqlx"
)

type sqlThumbnail struct {
	Id           string `db:"id"`
	BinaryFileId string `db:"binary_file_id"`
	Data         []byte `db:"data"`
	MimeType     string `db:"mime_type"`
}

type DefaultThumbnailRepo struct {
	db *sqlx.DB
}

func NewThumbnailRepo(db Database) repository.ThumbnailRepository {
	sqlxdb := db.Get()

	// Create table if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS thumbnails (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			binary_file_id VARCHAR(255) NOT NULL UNIQUE,
			data BLOB NOT NULL,
			mime_type VARCHAR(16) NOT NULL
		);
	`)

	return &DefaultThumbnailRepo{db: sqlxdb}
}

// Delete implements repository.ThumbnailRepository.
func (d *DefaultThumbnailRepo) Delete(binaryFileId string) error {
	_, err := d.db.Exec("DELETE FROM thumbnails WHERE binary_file_id = ?", binaryFileId)
	return err
}

// Get implements repository.ThumbnailRepository.
func (d *DefaultThumbnailRepo) Get(binaryFileId string) (*model.Thumbnail, error) {
	var sqlFile sqlThumbnail
	err := d.db.Get(&sqlFile, "SELECT * FROM thumbnails WHERE binary_file_id = ?", binaryFileId)
	if err != nil {
		return nil, err
	}
	return &model.Thumbnail{
		Id:           sqlFile.Id,
		BinaryFileId: sqlFile.BinaryFileId,
		Data:         sqlFile.Data,
		MimeType:     sqlFile.MimeType,
	}, nil
}

// Save implements repository.ThumbnailRepository.
func (d *DefaultThumbnailRepo) Save(binaryFileId string, mimeType string, data []byte) (*model.Thumbnail, error) {
	_, err := d.db.Exec("INSERT OR REPLACE INTO thumbnails (binary_file_id, data, mime_type) VALUES (?,?,?)", binaryFileId, data, mimeType)
	return &model.Thumbnail{
		BinaryFileId: binaryFileId,
		MimeType:     mimeType,
		Data:         data,
	}, err
}
