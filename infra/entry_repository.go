package infra

import (
	"encoding/json"
	"errors"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqlEntry struct {
	Id          string     `db:"id"`
	Type        string     `db:"type"`
	Content     string     `db:"content"`
	PublishedAt *time.Time `db:"published_at"`
	MetaData    *string    `db:"meta_data"`
}

type DefaultEntryRepo struct {
	types map[string]model.Entry
	db    *sqlx.DB
}

// Create implements repository.EntryRepository.
func (r *DefaultEntryRepo) Create(entry model.Entry) error {
	exEntry, _ := r.FindById(entry.ID())
	if exEntry != nil {
		return errors.New("entry already exists")
	}

	t := r.entryType(entry)
	if _, ok := r.types[t]; !ok {
		return errors.New("entry type not registered")
	}

	var metaDataJson []byte
	if entry.MetaData() != nil {
		metaDataJson, _ = json.Marshal(entry.MetaData())
	}

	_, err := r.db.Exec("INSERT INTO entries (id, type, content, published_at, meta_data) VALUES (?, ?, ?, ?, ?)", entry.ID(), t, entry.Content(), entry.PublishedAt(), metaDataJson)
	return err
}

// Delete implements repository.EntryRepository.
func (r *DefaultEntryRepo) Delete(entry model.Entry) error {
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

// FindById implements repository.EntryRepository.
func (r *DefaultEntryRepo) FindById(id string) (model.Entry, error) {
	data := sqlEntry{}
	err := r.db.Get(&data, "SELECT * FROM entries WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	if data.Id == "" {
		return nil, nil
	}
	return r.sqlEntryToEntry(data)
}

// RegisterEntryType implements repository.EntryRepository.
func (r *DefaultEntryRepo) RegisterEntryType(entry model.Entry) error {
	t := r.entryType(entry)
	if _, ok := r.types[t]; ok {
		return errors.New("entry type already registered")
	}
	r.types[t] = entry
	return nil
}

// Update implements repository.EntryRepository.
func (r *DefaultEntryRepo) Update(entry model.Entry) error {
	exEntry, _ := r.FindById(entry.ID())
	if exEntry == nil {
		return errors.New("entry not found")
	}

	t := r.entryType(entry)
	if _, ok := r.types[t]; !ok {
		return errors.New("entry type not registered")
	}

	var metaDataJson []byte
	if entry.MetaData() != nil {
		metaDataJson, _ = json.Marshal(entry.MetaData())
	}

	_, err := r.db.Exec("UPDATE entries SET content = ?, published_at = ?, meta_data = ? WHERE id = ?", entry.Content(), entry.PublishedAt(), metaDataJson, entry.ID())
	return err
}

func NewEntryRepository(db Database) repository.EntryRepository {
	sqlxdb := db.Get()

	// Create tables if not exists
	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS entries (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			content TEXT NOT NULL,
			published_at DATETIME,
			meta_data TEXT NOT NULL
		);
	`)

	return &DefaultEntryRepo{
		types: map[string]model.Entry{},
		db:    sqlxdb,
	}
}

func (r *DefaultEntryRepo) entryType(entry model.Entry) string {
	return reflect.TypeOf(entry).Elem().Name()
}

func (r *DefaultEntryRepo) sqlEntryToEntry(entry sqlEntry) (model.Entry, error) {
	e, ok := r.types[entry.Type]
	if !ok {
		return nil, errors.New("entry type not registered")
	}
	metaData := reflect.New(reflect.TypeOf(e.MetaData()).Elem()).Interface()
	json.Unmarshal([]byte(*entry.MetaData), metaData)
	e.Create(entry.Id, entry.Content, entry.PublishedAt, metaData)
	return e, nil
}
