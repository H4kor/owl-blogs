package infra

import (
	"encoding/json"
	"errors"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sqlInteraction struct {
	Id        string    `db:"id"`
	Type      string    `db:"type"`
	EntryId   string    `db:"entry_id"`
	CreatedAt time.Time `db:"created_at"`
	MetaData  *string   `db:"meta_data"`
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
func (repo *DefaultInteractionRepo) Create(interaction model.Interaction) error {
	t, err := repo.typeRegistry.TypeName(interaction)
	if err != nil {
		return errors.New("interaction type not registered")
	}

	if interaction.ID() == "" {
		interaction.SetID(uuid.New().String())
	}

	var metaDataJson []byte
	if interaction.MetaData() != nil {
		metaDataJson, _ = json.Marshal(interaction.MetaData())
	}
	metaDataStr := string(metaDataJson)

	_, err = repo.db.NamedExec(`
		INSERT INTO interactions (id, type, entry_id, created_at, meta_data)
		VALUES (:id, :type, :entry_id, :created_at, :meta_data)
	`, sqlInteraction{
		Id:        interaction.ID(),
		Type:      t,
		EntryId:   interaction.EntryID(),
		CreatedAt: interaction.CreatedAt(),
		MetaData:  &metaDataStr,
	})

	return err
}

// Delete implements repository.InteractionRepository.
func (*DefaultInteractionRepo) Delete(interaction model.Interaction) error {
	panic("unimplemented")
}

// FindAll implements repository.InteractionRepository.
func (repo *DefaultInteractionRepo) FindAll(entryId string) ([]model.Interaction, error) {
	data := []sqlInteraction{}
	err := repo.db.Select(&data, "SELECT * FROM interactions WHERE entry_id = ?", entryId)
	if err != nil {
		return nil, err
	}

	interactions := []model.Interaction{}
	for _, d := range data {
		i, err := repo.sqlInteractionToInteraction(d)
		if err != nil {
			return nil, err
		}
		interactions = append(interactions, i)
	}

	return interactions, nil
}

// FindById implements repository.InteractionRepository.
func (repo *DefaultInteractionRepo) FindById(id string) (model.Interaction, error) {
	data := sqlInteraction{}
	err := repo.db.Get(&data, "SELECT * FROM interactions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	if data.Id == "" {
		return nil, errors.New("interaction not found")
	}
	return repo.sqlInteractionToInteraction(data)
}

// Update implements repository.InteractionRepository.
func (repo *DefaultInteractionRepo) Update(interaction model.Interaction) error {
	exInter, _ := repo.FindById(interaction.ID())
	if exInter == nil {
		return errors.New("interaction not found")
	}

	_, err := repo.typeRegistry.TypeName(interaction)
	if err != nil {
		return errors.New("interaction type not registered")
	}

	var metaDataJson []byte
	if interaction.MetaData() != nil {
		metaDataJson, _ = json.Marshal(interaction.MetaData())
	}

	_, err = repo.db.Exec("UPDATE interactions SET entry_id = ?, meta_data = ? WHERE id = ?", interaction.EntryID(), metaDataJson, interaction.ID())

	return err
}

func (repo *DefaultInteractionRepo) sqlInteractionToInteraction(interaction sqlInteraction) (model.Interaction, error) {
	i, err := repo.typeRegistry.Type(interaction.Type)
	if err != nil {
		return nil, errors.New("interaction type not registered")
	}
	metaData := reflect.New(reflect.TypeOf(i.MetaData()).Elem()).Interface()
	json.Unmarshal([]byte(*interaction.MetaData), metaData)
	i.SetID(interaction.Id)
	i.SetEntryID(interaction.EntryId)
	i.SetCreatedAt(interaction.CreatedAt)
	i.SetMetaData(metaData)

	return i, nil

}
