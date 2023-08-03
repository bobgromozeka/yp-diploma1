package app

import (
	"database/sql"

	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

type App struct {
	Storage storage.Storage
	DB      *sql.DB
}

func New(s storage.Storage, db *sql.DB) App {
	return App{
		s,
		db,
	}
}
