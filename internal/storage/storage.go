package storage

import (
	"context"
	"errors"

	"github.com/bobgromozeka/yp-diploma1/internal/models"
)

var (
	ErrUserAlreadyExists   = errors.New("users already exists")
	ErrUserNotFound        = errors.New("users not found")
	ErrOrderNotFound       = errors.New("order not found")
	ErrOrderAlreadyCreated = errors.New("order already created")
	ErrOrderForeign        = errors.New("order foreign")
	ErrInsufficientFunds   = errors.New("insufficient funds")
)

type Storage interface {
	CreateUser(ctx context.Context, login string, password string) error
	AuthUser(ctx context.Context, login string, password string) (int64, error)
	CreateOrder(ctx context.Context, number string, userID int64) error
	GetUserOrders(ctx context.Context, userID int64) ([]models.Order, error)
	GetLatestUnprocessedOrders(ctx context.Context, count int) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status string, accrual *float64) error
	Withdraw(ctx context.Context, userID int64, orderNumber string, sum int) error
	GetUserWithdrawalsSum(ctx context.Context, userID int64) (int, error)
	GetUserBalance(ctx context.Context, userID int64) (int, error)
	GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdrawal, error)
}
