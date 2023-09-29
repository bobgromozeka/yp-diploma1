package testutils

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
)

func MatchContext() gomock.Matcher {
	ctx := reflect.TypeOf((*context.Context)(nil)).Elem()
	return gomock.AssignableToTypeOf(ctx)
}
