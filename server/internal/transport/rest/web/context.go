package web

import (
	"context"
	stderr "errors"

	"github.com/skrpld/NearBeee/internal/core/models/entities"
	"github.com/skrpld/NearBeee/pkg/errors"
)

type ctxKey int

const (
	CtxUserKey ctxKey = iota
	CtxErrorKey
)

func GetHttpErrorFromCtx(ctx context.Context) *errors.HttpError {
	var val *errors.HttpError
	_ = stderr.As(ctx.Value(CtxErrorKey).(*errors.HttpError), &val)
	return val
}

func GetUserFromCtx(ctx context.Context) (*entities.User, error) {
	ctxUser := ctx.Value(CtxUserKey)
	user, ok := ctxUser.(*entities.User)
	if !ok {
		return nil, errors.ErrNoPermissions
	}
	return user, nil
}
