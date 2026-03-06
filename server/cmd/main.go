package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skrpld/NearBeee/internal/core/database/mongodb"
	"github.com/skrpld/NearBeee/internal/core/database/postgres"
	"github.com/skrpld/NearBeee/internal/core/logger"
	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/transport/rest/servers"

	"github.com/skrpld/NearBeee/internal/config"
)

// TODO:
//  redis
//  kafka?
//  messages с mongo начать делать
//  создать воркеров и пул работ для распараллеливания
//  проверка широты\долготы

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
		zapLogger.Error("mongodb.NewMongoDB", logger.Error(err))
		return
	}

	postgresDB, err := postgres.NewPostgresDB(cfg.PostgresConfig, dbCtx)
	if err != nil {
		zapLogger.Error("postgres.NewPostgresDB", logger.Error(err))
		return
	}

	repo := repository.NewNearBeeeRepository(postgresDB, mongoDB)

	server, err := servers.NewHttpServer(cfg.HttpServerConfig, repo, zapLogger)
	if err != nil {
		zapLogger.Error("servers.NewNearBeeeServer", logger.Error(err))
		return
	}

	graceChan := make(chan os.Signal, 1)
	signal.Notify(graceChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			zapLogger.Error("server.Start", logger.Error(err))
		}
	}()
	<-graceChan

	if err = server.Stop(); err != nil {
		zapLogger.Error("server.Stop", logger.Error(err))
	}

	if err = postgresDB.Close(); err != nil {
		zapLogger.Error("postgresDB.Close", logger.Error(err))
	} else {
		zapLogger.With(logger.Time("stopped_at", time.Now())).Info("postgresDB closed")
	}
	if err = mongoDB.Close(); err != nil {
		zapLogger.Error("mongoDB.Close", logger.Error(err))
	} else {
		zapLogger.With(logger.Time("stopped_at", time.Now())).Info("mongoDB closed")
	}
}
