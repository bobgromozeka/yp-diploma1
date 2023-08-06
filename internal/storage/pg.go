package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bobgromozeka/yp-diploma1/internal/hash"
	"github.com/bobgromozeka/yp-diploma1/internal/models"
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

func (s PgStorage) AuthUser(ctx context.Context, login string, password string) (int64, error) {
	hashedPwd := hash.Sha256([]byte(password))
	row := s.db.QueryRowContext(ctx, "select id from users where login = $1 and password = $2", login, hashedPwd)

	if row.Err() != nil {
		return 0, row.Err()
	}

	var ID int64

	if scanErr := row.Scan(&ID); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return 0, UserNotFound
		}
		return 0, scanErr
	}

	return ID, nil
}

func (s PgStorage) CreateOrder(ctx context.Context, number string, userID int64) error {
	existingOrder, existingOrderErr := s.GetOrder(ctx, number)
	if existingOrderErr != nil && !errors.Is(existingOrderErr, OrderNotFound) {
		return existingOrderErr
	} else if existingOrderErr == nil {
		if existingOrder.UserID == userID {
			return OrderAlreadyCreated
		}
		return OrderForeign
	}

	_, createErr := s.db.ExecContext(
		ctx, "insert into orders(user_id, number, status, created_at) values($1,$2,$3,$4)", userID, number,
		models.OrderFirstStatus, time.Now(),
	)
	if createErr != nil {
		return createErr
	}

	return nil
}

func (s PgStorage) GetOrder(ctx context.Context, number string) (models.Order, error) {
	var o models.Order

	row := s.db.QueryRowContext(
		ctx, "select id, user_id, number, status, created_at, updated_at from orders where number = $1", number,
	)

	if row.Err() != nil {
		return o, row.Err()
	}

	if scanErr := row.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.CreatedAt, &o.UpdatedAt); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return o, OrderNotFound
		}
		return o, scanErr
	}

	return o, nil
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
    			id bigserial primary key,
    			login varchar(255) unique,
    			password varchar(255)
			)`,
	)
	if usersTableError != nil {
		return usersTableError
	}

	_, ordersTableError := tx.ExecContext(
		ctx,
		`create table if not exists orders(
    			id bigserial primary key,
    			user_id bigint,
    			number varchar(255) NOT NULL,
    			status varchar(255) NOT NULL,
    			accrual int,
    			created_at timestamp NOT NULL,
    			updated_at timestamp,
    			constraint fk_user
            		foreign key (user_id)
                    references users(id)
			)`,
	)
	if ordersTableError != nil {
		return ordersTableError
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
