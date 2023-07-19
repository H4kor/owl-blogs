package infra

import (
	"owl-blogs/domain/model"

	"github.com/jmoiron/sqlx"
)

type sqlAuthor struct {
	Name         string `db:"name"`
	PasswordHash string `db:"password_hash"`
}

type DefaultAuthorRepo struct {
	db *sqlx.DB
}

func NewDefaultAuthorRepo(db Database) *DefaultAuthorRepo {
	sqlxdb := db.Get()

	// Create table if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS authors (
			name TEXT PRIMARY KEY,
			password_hash TEXT NOT NULL
		);
	`)

	return &DefaultAuthorRepo{
		db: sqlxdb,
	}
}

// FindByName implements repository.AuthorRepository.
func (r *DefaultAuthorRepo) FindByName(name string) (*model.Author, error) {
	var author sqlAuthor
	err := r.db.Get(&author, "SELECT * FROM authors WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	return &model.Author{
		Name:         author.Name,
		PasswordHash: author.PasswordHash,
	}, nil
}

// Create implements repository.AuthorRepository.
func (r *DefaultAuthorRepo) Create(name string, passwordHash string) (*model.Author, error) {
	author := sqlAuthor{
		Name:         name,
		PasswordHash: passwordHash,
	}
	_, err := r.db.NamedExec("INSERT INTO authors (name, password_hash) VALUES (:name, :password_hash)", author)
	if err != nil {
		return nil, err
	}
	return &model.Author{
		Name:         author.Name,
		PasswordHash: author.PasswordHash,
	}, nil
}
