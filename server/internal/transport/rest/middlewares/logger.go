package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/skrpld/NearBeee/internal/core/logger"
	"github.com/skrpld/NearBeee/internal/transport/rest/web"
	"github.com/skrpld/NearBeee/pkg/errors"
)

//Потом стоит переделать

func LoggerMiddleware(zapLogger logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqLogger := zapLogger.With(
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("request_id", r.Header.Get("X-Request-ID")))

			reqLogger.Info("request started")

			httpError := &errors.HttpError{}
			ctx := context.WithValue(r.Context(), web.CtxErrorKey, httpError)

			next.ServeHTTP(w, r.WithContext(ctx))

			duration := time.Since(start)

			logFields := []logger.Field{
				logger.Duration("duration", duration),
				logger.String("duration_human", duration.String()),
			}

			if httpError.HasError() {
				logFields = append(logFields, logger.String("error", httpError.Error()), logger.Int("status_code", httpError.Code))

				if httpError.Code == http.StatusInternalServerError {
					httpError.Err = errors.ErrInternalServer
				}

				w.WriteHeader(httpError.Code)
				payload := errors.MarshalError(httpError)

				if _, err := w.Write(payload); err != nil {
					logFields = append(logFields, logger.String("write_error", err.Error()))

					http.Error(w, errors.ErrInternalServer.Error(), http.StatusInternalServerError)
				}

				reqLogger.With(logFields...).Error("request failed")
				return
			}

			reqLogger.With(logFields...).Info("request completed")
		})
	}
}
