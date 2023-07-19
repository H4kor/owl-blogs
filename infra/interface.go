package infra

import "github.com/jmoiron/sqlx"

type Database interface {
	Get() *sqlx.DB
}
