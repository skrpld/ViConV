package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HttpError struct {
	Err  error `json:"-"`
	Code int   `json:"-"`
}

func NewHttpError(err error, status int) error {
	return &HttpError{err, status}
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ErrUnknownError.Error()
}

func (e *HttpError) HasError() bool {
	return e != nil && e.Err != nil && e.Code != 0
}

func ParseHttpError(err error) *HttpError {
	var val *HttpError
	if errors.As(err, &val) {
		return val
	}
	return &HttpError{err, http.StatusInternalServerError}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (e *HttpError) MarshalJSON() ([]byte, error) {
	resp := ErrorResponse{ErrUnknownError.Error()}

	if e.HasError() {
		resp.Error = e.Error()
	}
	return json.Marshal(resp)
}

func MarshalError(err error) []byte {
	httpError := ParseHttpError(err)
	data, _ := json.Marshal(httpError)
	return data
}

var (
	ErrInvalidEmail                = NewHttpError(errors.New("invalid email"), http.StatusBadRequest)
	ErrUserAlreadyExists           = NewHttpError(errors.New("user already exists"), http.StatusBadRequest)
	ErrInvalidPassword             = NewHttpError(errors.New("invalid password"), http.StatusBadRequest)
	ErrInternalServer              = NewHttpError(errors.New("internal server error"), http.StatusInternalServerError)
	ErrInvalidToken                = NewHttpError(errors.New("invalid token"), http.StatusUnauthorized)
	ErrExpiredToken                = NewHttpError(errors.New("expired token"), http.StatusUnauthorized)
	ErrIdempotencyKeyAlreadyExists = NewHttpError(errors.New("idempotency key already exists"), http.StatusBadRequest)
	ErrInvalidPostId               = NewHttpError(errors.New("invalid post id"), http.StatusBadRequest)
	ErrNoPermissions               = NewHttpError(errors.New("no permissions"), http.StatusForbidden)
	ErrInvalidPostType             = NewHttpError(errors.New("invalid post type"), http.StatusBadRequest)
	ErrInvalidCoords               = NewHttpError(errors.New("invalid coordinates"), http.StatusBadRequest)
	ErrUnknownError                = NewHttpError(errors.New("unknown error"), http.StatusInternalServerError)
)
