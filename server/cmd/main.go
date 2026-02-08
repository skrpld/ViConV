package main

import (
	"context"
	"viconv/internal/config"
	"viconv/internal/database/mongodb"
	"viconv/internal/logger"
	"viconv/internal/transport/grpc/servers"
)

// TODO: real-time config
func main() {
	ctx := context.Background()

	zapLogger := logger.NewLogger()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := mongodb.NewDatabase(cfg.MongodbConfig, ctx)
	if err != nil {
		panic(err)
	}

	grpcServer, err := servers.NewViconvServer(cfg.ViconvServerConfig, &ctx, db, zapLogger)
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Start(); err != nil {
		panic(err)
	}
}
