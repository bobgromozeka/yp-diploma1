package app

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

type App struct {
	Storage storage.Storage
	DB      *sql.DB
	Logger  *zap.SugaredLogger
}

func New(s storage.Storage, db *sql.DB, logger *zap.SugaredLogger) App {
	return App{
		s,
		db,
		logger,
	}
}
