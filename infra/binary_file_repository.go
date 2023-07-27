package infra

import (
	"fmt"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type sqlBinaryFile struct {
	Id      string  `db:"id"`
	Name    string  `db:"name"`
	EntryId *string `db:"entry_id"`
	Data    []byte  `db:"data"`
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
			entry_id VARCHAR(255),
			data BLOB NOT NULL
		);
	`)

	return &DefaultBinaryFileRepo{db: sqlxdb}
}

// Create implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) Create(name string, data []byte, entry model.Entry) (*model.BinaryFile, error) {
	parts := strings.Split(name, ".")
	fileName := strings.Join(parts[:len(parts)-1], ".")
	fileExt := parts[len(parts)-1]
	id := fileName + "." + fileExt

	// check if id exists
	var count int
	err := repo.db.Get(&count, "SELECT COUNT(*) FROM binary_files WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		counter := 1
		for {
			id = fmt.Sprintf("%s-%d.%s", fileName, counter, fileExt)
			err := repo.db.Get(&count, "SELECT COUNT(*) FROM binary_files WHERE id = ?", id)
			if err != nil {
				return nil, err
			}
			if count == 0 {
				break
			}
			counter++
		}
	}

	var entryId *string
	if entry != nil {
		eId := entry.ID()
		entryId = &eId
	}

	_, err = repo.db.Exec("INSERT INTO binary_files (id, name, entry_id, data) VALUES (?, ?, ?, ?)", id, name, entryId, data)
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

// FindByNameForEntry implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) FindByNameForEntry(name string, entry model.Entry) (*model.BinaryFile, error) {
	var sqlFile sqlBinaryFile
	err := repo.db.Get(&sqlFile, "SELECT * FROM binary_files WHERE name = ? AND entry_id = ?", name, entry.ID())
	if err != nil {
		return nil, err
	}
	return &model.BinaryFile{Id: sqlFile.Id, Name: sqlFile.Name, Data: sqlFile.Data}, nil
}

// ListIds implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) ListIds() ([]string, error) {
	var ids []string
	err := repo.db.Select(&ids, "SELECT id FROM binary_files")
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// Delete implements repository.BinaryRepository
func (repo *DefaultBinaryFileRepo) Delete(binary *model.BinaryFile) error {
	id := binary.Id
	println("Deleting binary file", id)
	_, err := repo.db.Exec("DELETE FROM binary_files WHERE id = ?", id)
	return err
}
