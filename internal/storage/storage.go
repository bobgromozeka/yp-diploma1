package storage

import (
	"context"
	"errors"
)

var (
	UserAlreadyExists   = errors.New("users already exists")
	UserNotFound        = errors.New("users not found")
	OrderNotFound       = errors.New("order not found")
	OrderAlreadyCreated = errors.New("order already created")
	OrderForeign        = errors.New("order foreign")
)

type Storage interface {
	CreateUser(ctx context.Context, login string, password string) error
	AuthUser(ctx context.Context, login string, password string) (int64, error)
	CreateOrder(ctx context.Context, number string, userID int64) error
}