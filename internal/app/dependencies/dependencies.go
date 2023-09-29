package dependencies

import (
	"database/sql"

	"go.uber.org/zap"

	"github.com/bobgromozeka/yp-diploma1/internal/storage"
)

type D struct {
	UsersStorage       storage.UsersStorage
	OrdersStorage      storage.OrdersStorage
	WithdrawalsStorage storage.WithdrawalsStorage
	DB                 *sql.DB
	Logger             *zap.SugaredLogger
}
