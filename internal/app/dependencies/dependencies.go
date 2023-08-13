package dependencies

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

type D struct {
	Storage storage.Storage
	DB      *sql.DB
	Logger  *zap.SugaredLogger
}

func New(s storage.Storage, db *sql.DB, logger *zap.SugaredLogger) D {
	return D{
		s,
		db,
		logger,
	}
}
