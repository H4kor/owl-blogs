package test

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type MockDb struct {
	db *sqlx.DB
}

func (d *MockDb) Get() *sqlx.DB {
	return d.db
}

func NewMockDb() *MockDb {
	db := sqlx.MustOpen("sqlite3", ":memory:")
	return &MockDb{db: db}
}
