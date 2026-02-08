package interceptors

import (
	"context"
	"strings"
	"viconv/internal/models/dto"
	"viconv/pkg/consts/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthService interface {
	AuthorizeUser(rows *dto.AuthorizeUserRequest) (*dto.AuthorizeUserResponse, error)
}

type AuthInterceptorHandler struct {
	srv AuthService
}

func NewAuthInterceptorHandler(srv AuthService) *AuthInterceptorHandler {
	return &AuthInterceptorHandler{
		srv: srv,
	}
}

var protectedMethods = []string{
	"LoginUser",
	"RegistrateUser",
	"RefreshUserToken",
} //TODO: refactor

func (h *AuthInterceptorHandler) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		isProtected := false
		for _, method := range protectedMethods {
			if strings.Contains(info.FullMethod, method) {
				isProtected = true
				break
			}
		}

		if !isProtected {

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, errors.ErrInvalidToken
			}

			header := md.Get("authorization")
			if len(header) == 0 {
				return nil, errors.ErrInvalidToken
			}

			token := strings.Split(header[0], " ")
			if len(token) != 2 {
				return nil, errors.ErrInvalidToken
			}

			if token[0] != "Bearer" {
				return nil, errors.ErrInvalidToken
			}

			user, err := h.srv.AuthorizeUser(&dto.AuthorizeUserRequest{
				AccessToken: token[1],
			})
			if err != nil {
				return nil, errors.ErrInvalidToken
			}

			ctx = context.WithValue(ctx, "user", user)
		}

		return handler(ctx, req)
	}
}
