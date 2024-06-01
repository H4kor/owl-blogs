package infra

import (
	"encoding/json"
	"owl-blogs/app"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqlEntry struct {
	Id          string     `db:"id"`
	Type        string     `db:"type"`
	PublishedAt *time.Time `db:"published_at"`
	MetaData    *string    `db:"meta_data"`
	AuthorId    string     `db:"author_id"`
}

type DefaultEntryRepo struct {
	typeRegistry *app.EntryTypeRegistry
	db           *sqlx.DB
}

func NewEntryRepository(db Database, register *app.EntryTypeRegistry) repository.EntryRepository {
	sqlxdb := db.Get()

	// Create tables if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			published_at DATETIME,
			author_id TEXT NOT NULL,
			meta_data TEXT NOT NULL
		);
	`)

	return &DefaultEntryRepo{
		db:           sqlxdb,
		typeRegistry: register,
	}
}

// Create implements repository.EntryRepository.
func (r *DefaultEntryRepo) Create(entry model.Entry) error {
	t, err := r.typeRegistry.TypeName(entry)
	if err != nil {
		return err
	}

	var metaDataJson []byte
	if entry.MetaData() != nil {
		metaDataJson, _ = json.Marshal(entry.MetaData())
	}

	if entry.ID() == "" {
		entry.SetID(uuid.New().String())
	}

	_, err = r.db.Exec("INSERT INTO entries (id, type, published_at, author_id, meta_data) VALUES (?, ?, ?, ?, ?)", entry.ID(), t, entry.PublishedAt(), entry.AuthorId(), metaDataJson)
	return err
}

// Delete implements repository.EntryRepository.
func (r *DefaultEntryRepo) Delete(entry model.Entry) error {
	if entry.ID() == "" {
		return app.ErrEntryNotFound
	}
	_, err := r.db.Exec("DELETE FROM entries WHERE id = ?", entry.ID())
	return err
}

// FindAll implements repository.EntryRepository.
func (r *DefaultEntryRepo) FindAll(types *[]string) ([]model.Entry, error) {
	filterStr := ""
	if types != nil {
		filters := []string{}
		for _, t := range *types {
			filters = append(filters, "type = '"+t+"'")
		}
		filterStr = strings.Join(filters, " OR ")
	}

	var entries []sqlEntry
	if filterStr != "" {
		err := r.db.Select(&entries, "SELECT * FROM entries WHERE "+filterStr)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&entries, "SELECT * FROM entries")
		if err != nil {
			return nil, err
		}
	}

	result := []model.Entry{}
	for _, entry := range entries {
		e, err := r.sqlEntryToEntry(entry)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, nil
}

// FindAll implements repository.EntryRepository.
func (r *DefaultEntryRepo) FindAllByTag(tag string) ([]model.Entry, error) {
	// tags are not persisted into database (yet)
	// therefore we have to get all entries and filter in application
	entries, err := r.FindAll(nil)
	if err != nil {
		return nil, err
	}
	tagEntries := make([]model.Entry, 0)
	for _, e := range entries {
		if slices.Contains(e.Tags(), tag) {
			tagEntries = append(tagEntries, e)
		}
	}
	return tagEntries, nil
}

// FindById implements repository.EntryRepository.
func (r *DefaultEntryRepo) FindById(id string) (model.Entry, error) {
	data := sqlEntry{}
	err := r.db.Get(&data, "SELECT * FROM entries WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	if data.Id == "" {
		return nil, app.ErrEntryNotFound
	}
	return r.sqlEntryToEntry(data)
}

// Update implements repository.EntryRepository.
func (r *DefaultEntryRepo) Update(entry model.Entry) error {
	exEntry, _ := r.FindById(entry.ID())
	if exEntry == nil {
		return app.ErrEntryNotFound
	}

	_, err := r.typeRegistry.TypeName(entry)
	if err != nil {
		return err
	}

	var metaDataJson []byte
	if entry.MetaData() != nil {
		metaDataJson, _ = json.Marshal(entry.MetaData())
	}

	_, err = r.db.Exec("UPDATE entries SET published_at = ?, author_id = ?, meta_data = ? WHERE id = ?", entry.PublishedAt(), entry.AuthorId(), metaDataJson, entry.ID())
	return err
}

func (r *DefaultEntryRepo) sqlEntryToEntry(entry sqlEntry) (model.Entry, error) {
	e, err := r.typeRegistry.Type(entry.Type)
	if err != nil {
		return nil, err
	}
	metaData := reflect.New(reflect.TypeOf(e.MetaData()).Elem()).Interface().(model.EntryMetaData)
	json.Unmarshal([]byte(*entry.MetaData), metaData)
	e.SetID(entry.Id)
	e.SetPublishedAt(entry.PublishedAt)
	e.SetMetaData(metaData)
	e.SetAuthorId(entry.AuthorId)
	return e, nil
}
