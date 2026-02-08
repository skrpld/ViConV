package main

import (
	"context"
	"time"
	"viconv/internal/config"
	"viconv/internal/database/mongodb"
	"viconv/internal/database/postgres"
	"viconv/internal/logger"
	"viconv/internal/transport/grpc/servers"

	"go.uber.org/zap"
)

// TODO: real-time config

func main() {
	ctx := context.Background()

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	zapLogger, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		zapLogger.Error("config.NewConfig", zap.Error(err))
		panic(err)
	}

	mongoDB, err := mongodb.NewMongoDB(cfg.MongoDBConfig, dbCtx)
	if err != nil {
		zapLogger.Error("mongodb.NewMongoDB", zap.Error(err))
		panic(err)
	}

	postgresDB, err := postgres.NewPostgresDB(cfg.PostgresConfig, dbCtx)
	if err != nil {
		zapLogger.Error("postgres.NewPostgresDB", zap.Error(err))
		panic(err)
	}

	grpcServer, err := servers.NewViconvServer(cfg.ViconvServerConfig, &ctx, mongoDB, postgresDB, zapLogger)
	if err != nil {
		zapLogger.Error("servers.NewViconvServer", zap.Error(err))
		panic(err)
	}
	if err := grpcServer.Start(); err != nil {
		panic(err)
	}
}
