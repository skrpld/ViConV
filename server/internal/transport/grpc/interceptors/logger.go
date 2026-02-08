package interceptors

import (
	"context"
	"fmt"
	"time"
	"viconv/internal/logger"
	"viconv/pkg/consts/errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor(zapLogger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		requestId := start.UnixNano()
		method := info.FullMethod

		reqLogger := zapLogger.With(
			zap.String("method", method),
			zap.String("request_id", fmt.Sprint(requestId)))
		reqLogger.Info("gRPC request started")

		err := grpc.SetHeader(ctx, metadata.Pairs("X-Request-ID", fmt.Sprint(requestId)))
		if err != nil {
			zap.Error(err)
		}

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		logFields := []zap.Field{
			zap.String("method", method),
			zap.Duration("duration", duration),
			zap.String("duration_formated", duration.String()),
		}

		if err != nil {
			logFields = append(logFields,
				zap.Error(err),
				zap.String("status_code", status.Code(err).String()),
			)
			reqLogger.With(logFields...).Error("gRPC request failed")

			_, ok := status.FromError(err)
			if !ok {
				err = errors.ErrInternalServer
			}

		} else {
			reqLogger.With(logFields...).Info("gRPC request completed")
		}
		return resp, err
	}
}
