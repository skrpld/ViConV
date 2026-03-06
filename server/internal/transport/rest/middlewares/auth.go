package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/skrpld/NearBeee/internal/core/models/dto"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
	"github.com/skrpld/NearBeee/pkg/errors"
)

type AuthService interface {
	AuthorizeUser(rows *dto.AuthorizeUserRequest) (*dto.AuthorizeUserResponse, error)
}
type AuthMiddlewareHandler struct {
	srv AuthService
}

func NewAuthMiddlewareHandler(srvAuth AuthService) *AuthMiddlewareHandler {
	return &AuthMiddlewareHandler{srvAuth}
}

func (a *AuthMiddlewareHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpError := web.GetHttpErrorFromCtx(r.Context())

		header := r.Header.Get("Authorization")
		tokenParts := strings.Split(header, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			parsedError := errors.ParseHttpError(errors.ErrInvalidToken)
			httpError.Err = parsedError.Err
			httpError.Code = parsedError.Code

			return
		}

		token := tokenParts[1]

		user, err := a.srv.AuthorizeUser(&dto.AuthorizeUserRequest{AccessToken: token})
		if err != nil {
			parsedError := errors.ParseHttpError(err)
			httpError.Err = parsedError.Err
			httpError.Code = parsedError.Code

			return
		}

		ctx := context.WithValue(r.Context(), web.CtxUserKey, user.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
