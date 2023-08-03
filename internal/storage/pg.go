package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bobgromozeka/yp-diploma1/internal/hash"
)

type PgStorage struct {
	db *sql.DB
}

func NewPGStorage(db *sql.DB) Storage {
	return PgStorage{
		db,
	}
}

func (s PgStorage) CreateUser(ctx context.Context, login string, password string) error {
	hashedPwd := hash.Sha256([]byte(password))
	_, err := s.db.ExecContext(ctx, "insert into users (login, password) values ($1, $2)", login, hashedPwd)

	if err != nil && IsExactCode(err, pgerrcode.UniqueViolation) {
		return UserAlreadyExists
	}

	return err
}

func (s PgStorage) AuthUser(ctx context.Context, login string, password string) (bool, error) {
	hashedPwd := hash.Sha256([]byte(password))
	row := s.db.QueryRowContext(ctx, "select * from users where login = $1 and password = $2", login, hashedPwd)

	if row.Err() != nil && row.Err() != sql.ErrNoRows {
		return false, row.Err()
	}

	if row.Err() == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func Bootstrap(db *sql.DB) error {
	ctx := context.Background()
	tx, txErr := db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	_, usersTableError := tx.ExecContext(
		ctx,
		`create table if not exists users(
    			login varchar(255),
    			password varchar(255),
    			primary key (login)
			)`,
	)
	if usersTableError != nil {
		return usersTableError
	}

	tx.Commit()

	return nil
}

func IsExactType(err error, errFunc func(string) bool) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && errFunc(pgErr.Code)
}

func IsExactCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && pgErr.Code == code
}
