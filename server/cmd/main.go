package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"viconv/internal/database/mongodb"
	"viconv/internal/database/postgres"
	"viconv/internal/transport/grpc/servers"

	"viconv/internal/config"
	"viconv/internal/logger"

	"go.uber.org/zap"
)

// TODO: real-time config +- Доделать
//  redis
//  kafka/rabbitmq?
//  post message + в controller формирование поста

func main() {
	ctx := context.Background()

	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second) //TODO: как то переделать
	defer cancel()

	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.GetConfig()

	zapLogger, err := logger.NewLogger(cfg.LoggerConfig)
	if err != nil {
		panic(err)
	}

	mongoDB, err := mongodb.NewMongoDB(cfg.MongoDBConfig, dbCtx)
	if err != nil {
		zapLogger.Error("mongodb.NewMongoDB", zap.Error(err))
		return
	}

	postgresDB, err := postgres.NewPostgresDB(cfg.PostgresConfig, dbCtx)
	if err != nil {
		zapLogger.Error("postgres.NewPostgresDB", zap.Error(err))
		return
	}

	grpcServer, err := servers.NewViconvServer(cfg.ViconvServerConfig, &ctx, mongoDB, postgresDB, zapLogger)
	if err != nil {
		zapLogger.Error("servers.NewViconvServer", zap.Error(err))
		return
	}

	graceChan := make(chan os.Signal, 1)
	signal.Notify(graceChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Start(); err != nil {
			zapLogger.Error("grpcServer.Start", zap.Error(err))
		}
	}()
	<-graceChan

	grpcServer.Stop()

	if err = postgresDB.Close(); err != nil {
		zapLogger.Error("postgresDB.Close", zap.Error(err))
	} else {
		zapLogger.With(zap.Time("stopped_at", time.Now())).Info("postgresDB closed")
	}
	if err = mongoDB.Close(); err != nil {
		zapLogger.Error("mongoDB.Close", zap.Error(err))
	} else {
		zapLogger.With(zap.Time("stopped_at", time.Now())).Info("mongoDB closed")
	}
}
