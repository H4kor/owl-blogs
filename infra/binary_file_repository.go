package infra

import (
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sqlBinaryFile struct {
	Id   string `db:"id"`
	Name string `db:"name"`
	Data []byte `db:"data"`
}

type DefaultBinaryFileRepo struct {
	db *sqlx.DB
}

// NewBinaryFileRepo creates a new binary file repository
// It creates the table if not exists
func NewBinaryFileRepo(db Database) repository.BinaryRepository {
	sqlxdb := db.Get()

	// Create table if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS binary_files (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			data BLOB NOT NULL
		);
	`)

	return &DefaultBinaryFileRepo{db: sqlxdb}
}

// Create implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) Create(name string, data []byte) (*model.BinaryFile, error) {
	id := uuid.New().String()
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		ext := parts[len(parts)-1]
		id = id + "." + ext
	}

	_, err := repo.db.Exec("INSERT INTO binary_files (id, name, data) VALUES (?, ?, ?)", id, name, data)
	if err != nil {
		return nil, err
	}
	return &model.BinaryFile{Id: id, Name: name, Data: data}, nil
}

// FindById implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) FindById(id string) (*model.BinaryFile, error) {
	var sqlFile sqlBinaryFile
	err := repo.db.Get(&sqlFile, "SELECT * FROM binary_files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &model.BinaryFile{Id: sqlFile.Id, Name: sqlFile.Name, Data: sqlFile.Data}, nil
}
