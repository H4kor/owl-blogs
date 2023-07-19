package infra

import (
	"encoding/json"
	"owl-blogs/app/repository"

	"github.com/jmoiron/sqlx"
)

type DefaultConfigRepo struct {
	db *sqlx.DB
}

func NewConfigRepo(db Database) repository.ConfigRepository {
	sqlxdb := db.Get()

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS site_config (
			name TEXT PRIMARY KEY,
			config TEXT
		);
	`)

	return &DefaultConfigRepo{
		db: sqlxdb,
	}
}

// Get implements repository.SiteConfigRepository.
func (r *DefaultConfigRepo) Get(name string, result interface{}) error {
	data := []byte{}
	err := r.db.Get(&data, "SELECT config FROM site_config WHERE name = ?", name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, result)
}

// Update implements repository.SiteConfigRepository.
func (r *DefaultConfigRepo) Update(name string, siteConfig interface{}) error {
	jsonData, err := json.Marshal(siteConfig)
	if err != nil {
		return err
	}
	res, err := r.db.Exec("UPDATE site_config SET config = ? WHERE name = ?", jsonData, name)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		_, err = r.db.Exec("INSERT INTO site_config (name, config) VALUES (?, ?)", name, jsonData)
	}
	return err
}
