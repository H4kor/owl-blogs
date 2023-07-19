package infra

import (
	"github.com/jmoiron/sqlx"
)

type SqliteDatabase struct {
	db *sqlx.DB
}

func NewSqliteDB(path string) Database {
	db := sqlx.MustOpen("sqlite3", path)
	return &SqliteDatabase{db: db}
}

func (d *SqliteDatabase) Get() *sqlx.DB {
	return d.db
}
