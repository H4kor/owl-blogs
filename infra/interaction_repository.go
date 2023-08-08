package infra

import (
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"

	"github.com/jmoiron/sqlx"
)

type sqlInteraction struct {
	Id        string `db:"id"`
	Type      string `db:"type"`
	EntryId   string `db:"entry_id"`
	CreatedAt string `db:"created_at"`
	MetaData  string `db:"meta_data"`
}

type DefaultInteractionRepo struct {
	typeRegistry *app.InteractionTypeRegistry
	db           *sqlx.DB
}

func NewInteractionRepo(db Database, register *app.InteractionTypeRegistry) repository.InteractionRepository {
	sqlxdb := db.Get()

	// Create tables if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS interactions (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			entry_id TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			meta_data TEXT NOT NULL
		);
	`)

	return &DefaultInteractionRepo{
		db:           sqlxdb,
		typeRegistry: register,
	}
}

// Create implements repository.InteractionRepository.
func (*DefaultInteractionRepo) Create(interaction model.Interaction) error {
	panic("unimplemented")
}

// Delete implements repository.InteractionRepository.
func (*DefaultInteractionRepo) Delete(interaction model.Interaction) error {
	panic("unimplemented")
}

// FindAll implements repository.InteractionRepository.
func (*DefaultInteractionRepo) FindAll(entryId string) ([]model.Interaction, error) {
	panic("unimplemented")
}

// FindById implements repository.InteractionRepository.
func (*DefaultInteractionRepo) FindById(id string) (model.Interaction, error) {
	panic("unimplemented")
}

// Update implements repository.InteractionRepository.
func (*DefaultInteractionRepo) Update(interaction model.Interaction) error {
	panic("unimplemented")
}
