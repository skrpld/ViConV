package web

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/skrpld/NearBeee/pkg/errors"
)

type Handler func(r *http.Request) (any, error)

func Handle(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpError := GetHttpErrorFromCtx(r.Context())

		data, err := handler(r)

		if err != nil {
			parsedErr := errors.ParseHttpError(err)
			httpError.Err = parsedErr.Err
			httpError.Code = parsedErr.Code

			return
		}

		accessToken, ok := hasAccessToken(data)
		if ok {
			w.Header().Set("Authorization", "Bearer "+accessToken)
		}

		if data != nil {
			json.NewEncoder(w).Encode(data)
		}
	}
}

func hasAccessToken(v any) (string, bool) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return "", false
	}

	field := val.FieldByName("AccessToken")
	if !field.IsValid() {
		return "", false
	}

	if field.Kind() != reflect.String {
		return "", false
	}
	return field.String(), true
}
