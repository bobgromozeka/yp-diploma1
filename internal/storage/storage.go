package storage

import (
	"context"
	"errors"
)

var UserAlreadyExists = errors.New("user already exists")

type Storage interface {
	CreateUser(context.Context, string, string) error
	AuthUser(context.Context, string, string) (bool, error)
}
