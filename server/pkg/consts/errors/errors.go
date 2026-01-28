package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidEmail      = status.Error(codes.InvalidArgument, "error: invalid email")
	ErrUserAlreadyExists = status.Error(codes.AlreadyExists, "error: user already exists")
	ErrInvalidPassword   = status.Error(codes.InvalidArgument, "error: invalid password")
	ErrInternalServer    = status.Error(codes.Internal, "error: internal server error")
	ErrInvalidToken      = status.Error(codes.PermissionDenied, "error: invalid token")
	ErrExpiredToken      = status.Error(codes.PermissionDenied, "error: expired token")
)
