package infra

import (
	"encoding/json"
	"owl-blogs/app/repository"
	"owl-blogs/domain/model"

	"github.com/jmoiron/sqlx"
)

type DefaultSiteConfigRepo struct {
	db *sqlx.DB
}

func NewSiteConfigRepo(db Database) repository.SiteConfigRepository {
	sqlxdb := db.Get()

	sqlxdb.MustExec(`
		CREATE TABLE IF NOT EXISTS site_config (
			config TEXT
		);
	`)

	return &DefaultSiteConfigRepo{
		db: sqlxdb,
	}
}

// Get implements repository.SiteConfigRepository.
func (r *DefaultSiteConfigRepo) Get() (model.SiteConfig, error) {
	data := []byte{}
	err := r.db.Get(&data, "SELECT config FROM site_config LIMIT 1")
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return model.SiteConfig{}, nil
		}
		return model.SiteConfig{}, err
	}
	if len(data) == 0 {
		return model.SiteConfig{}, nil
	}
	config := model.SiteConfig{}
	err = json.Unmarshal(data, &config)
	return config, err
}

// Update implements repository.SiteConfigRepository.
func (r *DefaultSiteConfigRepo) Update(siteConfig model.SiteConfig) error {
	jsonData, err := json.Marshal(siteConfig)
	if err != nil {
		return err
	}
	res, err := r.db.Exec("UPDATE site_config SET config = ?", jsonData)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if rows == 0 {
		_, err = r.db.Exec("INSERT INTO site_config (config) VALUES (?)", jsonData)
	}
	return err
}
