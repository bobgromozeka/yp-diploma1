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
	tx, txErr := s.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	hashedPwd := hash.Sha256([]byte(password))
	row := tx.QueryRowContext(
		ctx, "insert into users (login, password) values ($1, $2) returning id", login, hashedPwd,
	)

	if row.Err() != nil && IsExactCode(row.Err(), pgerrcode.UniqueViolation) {
		return ErrUserAlreadyExists
	}

	var userLastInsertedID int64
	uliErr := row.Scan(&userLastInsertedID)
	if uliErr != nil {
		return uliErr
	}

	_, balanceErr := tx.ExecContext(ctx, "insert into user_balances (user_id) values ($1)", userLastInsertedID)
	if balanceErr != nil {
		return balanceErr
	}

	tx.Commit()
	return nil
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
			return 0, ErrUserNotFound
		}
		return 0, scanErr
	}

	return ID, nil
}

func (s PgStorage) CreateOrder(ctx context.Context, number string, userID int64) error {
	existingOrder, existingOrderErr := s.GetOrder(ctx, number)
	if existingOrderErr != nil && !errors.Is(existingOrderErr, ErrOrderNotFound) {
		return existingOrderErr
	} else if existingOrderErr == nil {
		if existingOrder.UserID == userID {
			return ErrOrderAlreadyCreated
		}
		return ErrOrderForeign
	}

	_, createErr := s.db.ExecContext(
		ctx, "insert into orders(user_id, number, status, uploaded_at) values($1,$2,$3,$4)", userID, number,
		models.OrderFirstStatus, time.Now(),
	)
	if createErr != nil {
		return createErr
	}

	return nil
}

func (s PgStorage) GetUserOrders(ctx context.Context, userID int64) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	rows, rowsErr := s.db.QueryContext(
		ctx, "select id, user_id, number, status, uploaded_at, updated_at from orders where user_id = $1", userID,
	)
	if rowsErr != nil {
		return orders, rowsErr
	}
	if rows.Err() != nil {
		return orders, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		if scanErr := rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.UploadedAt, &o.UpdatedAt); scanErr != nil {
			return orders, scanErr
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (s PgStorage) GetOrder(ctx context.Context, number string) (models.Order, error) {
	var o models.Order

	row := s.db.QueryRowContext(
		ctx, "select id, user_id, number, status, uploaded_at, updated_at from orders where number = $1", number,
	)

	if row.Err() != nil {
		return o, row.Err()
	}

	if scanErr := row.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.UploadedAt, &o.UpdatedAt); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return o, ErrOrderNotFound
		}
		return o, scanErr
	}

	return o, nil
}

func (s PgStorage) GetLatestUnprocessedOrders(ctx context.Context, count int) ([]models.Order, error) {
	orders := make([]models.Order, 0)
	rows, rowsErr := s.db.QueryContext(
		ctx,
		`select id, user_id, number, status, uploaded_at, updated_at from orders where status in ($1, $2)
                order by updated_at desc nulls first limit $3`, models.OrderStatusNew, models.OrderStatusProcessing,
		count,
	)
	if rowsErr != nil {
		return orders, rowsErr
	}
	if rows.Err() != nil {
		return orders, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		if scanErr := rows.Scan(&o.ID, &o.UserID, &o.Number, &o.Status, &o.UploadedAt, &o.UpdatedAt); scanErr != nil {
			return orders, scanErr
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (s PgStorage) UpdateOrderStatus(ctx context.Context, number string, status string, accrual *float64) error {
	tx, txErr := s.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(
		ctx, "update orders set status = $1, accrual = $2, updated_at = $3 where number = $4 returning user_id", status,
		accrual, time.Now(), number,
	)

	var userID int64

	if scanErr := row.Scan(&userID); scanErr != nil {
		return scanErr
	}
	if accrual != nil {
		_, balanceErr := tx.ExecContext(
			ctx, "update user_balances set balance = balance + $1 where user_id = $2", accrual, userID,
		)

		if balanceErr != nil {
			return balanceErr
		}
	}
	tx.Commit()

	return nil
}

func (s PgStorage) Withdraw(ctx context.Context, userID int64, orderNumber string, sum int) error {
	balanceRow := s.db.QueryRowContext(ctx, "select balance from user_balances where user_id = $1", userID)

	var balance int

	balanceErr := balanceRow.Scan(&balance)
	if balanceErr != nil {
		return balanceErr
	}

	if balance < sum {
		return ErrInsufficientFunds
	}

	tx, txErr := s.db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	_, withdrawErr := tx.ExecContext(
		ctx, "insert into withdrawals(user_id, order_number, sum, processed_at) values($1,$2,$3, $4)", userID,
		orderNumber, sum, time.Now(),
	)
	if withdrawErr != nil {
		return withdrawErr
	}

	_, updateBalanceErr := tx.ExecContext(
		ctx, "update user_balances set balance = balance - $1 where user_id = $2", sum, userID,
	)
	if updateBalanceErr != nil {
		return updateBalanceErr
	}
	tx.Commit()

	return nil
}

func (s PgStorage) GetUserWithdrawalsSum(ctx context.Context, userID int64) (int, error) {
	sumRow := s.db.QueryRowContext(ctx, "select sum(sum) from withdrawals where user_id = $1", userID)

	var sum int

	if scanErr := sumRow.Scan(&sum); scanErr != nil {
		return 0, scanErr
	}

	return sum, nil
}

func (s PgStorage) GetUserBalance(ctx context.Context, userID int64) (int, error) {
	balanceRow := s.db.QueryRowContext(ctx, "select balance from user_balances where user_id = $1", userID)

	var balance int

	if scanErr := balanceRow.Scan(&balance); scanErr != nil {
		return 0, scanErr
	}

	return balance, nil
}

func (s PgStorage) GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	withdrawalRows, withdrawalsErr := s.db.QueryContext(
		ctx, "select id, user_id, order_number, sum, processed_at from withdrawals where user_id = $1", userID,
	)
	if withdrawalsErr != nil {
		return withdrawals, withdrawalsErr
	}
	if withdrawalRows.Err() != nil {
		return withdrawals, withdrawalRows.Err()
	}
	defer withdrawalRows.Close()

	for withdrawalRows.Next() {
		var w models.Withdrawal
		if scanErr := withdrawalRows.Scan(&w.ID, &w.UserID, &w.OrderNumber, &w.Sum, &w.ProcessedAt); scanErr != nil {
			return withdrawals, scanErr
		}
		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}

func IsExactType(err error, errFunc func(string) bool) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && errFunc(pgErr.Code)
}

func IsExactCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return err != nil && errors.As(err, &pgErr) && pgErr.Code == code
}

func Bootstrap(db *sql.DB) error {
	ctx := context.Background()
	tx, txErr := db.BeginTx(ctx, nil)
	if txErr != nil {
		return txErr
	}
	defer tx.Rollback()

	usersTableError := createUsersTable(ctx, tx)
	if usersTableError != nil {
		return usersTableError
	}

	ordersTableError := createOrdersTable(ctx, tx)
	if ordersTableError != nil {
		return ordersTableError
	}

	userBalancesTableError := createUserBalancesTable(ctx, tx)
	if userBalancesTableError != nil {
		return userBalancesTableError
	}

	withdrawalsTableError := createWithdrawalsTable(ctx, tx)
	if withdrawalsTableError != nil {
		return withdrawalsTableError
	}

	tx.Commit()

	return nil
}

func createUsersTable(ctx context.Context, tx *sql.Tx) error {
	_, usersTableError := tx.ExecContext(
		ctx,
		`create table if not exists users(
    			id bigserial primary key,
    			login varchar(255) unique,
    			password varchar(255)
			)`,
	)

	return usersTableError
}

func createOrdersTable(ctx context.Context, tx *sql.Tx) error {
	_, ordersTableError := tx.ExecContext(
		ctx,
		`create table if not exists orders(
    			id bigserial primary key,
    			user_id bigint,
    			number varchar(255) NOT NULL,
    			status varchar(255) NOT NULL,
    			accrual int,
    			uploaded_at timestamp NOT NULL,
    			updated_at timestamp,
    			constraint fk_user
            		foreign key (user_id)
                    references users(id)
			)`,
	)

	return ordersTableError
}

func createUserBalancesTable(ctx context.Context, tx *sql.Tx) error {
	_, userBalancesError := tx.ExecContext(
		ctx,
		`create table if not exists user_balances(
    			id bigserial primary key,
    			user_id bigint,
    			balance int NOT NULL default 0,
    			constraint fk_user
            		foreign key (user_id)
                    references users(id)
			)`,
	)

	return userBalancesError
}

func createWithdrawalsTable(ctx context.Context, tx *sql.Tx) error {
	_, withdrawalsError := tx.ExecContext(
		ctx,
		`create table if not exists withdrawals(
    			id bigserial primary key,
    			user_id bigint,
    			order_number varchar(255) NOT NULL,
    			sum int NOT NULL,
    			processed_at timestamp NOT NULL,
    			constraint fk_user
            		foreign key (user_id)
                    references users(id)
			)`,
	)

	return withdrawalsError
}
