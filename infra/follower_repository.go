package infra

import (
	"owl-blogs/app/repository"

	"github.com/jmoiron/sqlx"
)

type sqlFollower struct {
	Follwer string `db:"follower"`
}

type DefaultFollowerRepo struct {
	db *sqlx.DB
}

func NewFollowerRepository(db Database) repository.FollowerRepository {
	sqlxdb := db.Get()

	// Create tables if not exists
	sqlxdb.MustExec(`
	CREATE TABLE IF NOT EXISTS followers (
		follower TEXT PRIMARY KEY
	);
	`)

	return &DefaultFollowerRepo{
		db: sqlxdb,
	}
}

// Add implements repository.FollowerRepository.
func (d *DefaultFollowerRepo) Add(follower string) error {
	_, err := d.db.Exec("INSERT INTO followers (follower) VALUES (?) ON CONFLICT DO NOTHING", follower)
	return err
}

// Remove implements repository.FollowerRepository.
func (d *DefaultFollowerRepo) Remove(follower string) error {
	_, err := d.db.Exec("DELETE FROM followers WHERE follower = ?", follower)
	return err
}

// All implements repository.FollowerRepository.
func (d *DefaultFollowerRepo) All() ([]string, error) {
	var followers []sqlFollower
	err := d.db.Select(&followers, "SELECT * FROM followers")
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, follower := range followers {
		result = append(result, follower.Follwer)
	}
	return result, nil
}
