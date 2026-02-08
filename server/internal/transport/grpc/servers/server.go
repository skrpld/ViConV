package servers

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/skrpld/NearBeee/internal/database/mongodb"
	"github.com/skrpld/NearBeee/internal/database/postgres"
	"github.com/skrpld/NearBeee/internal/logger"
	"github.com/skrpld/NearBeee/internal/repository"
	"github.com/skrpld/NearBeee/internal/service"
	"github.com/skrpld/NearBeee/internal/transport/grpc/controllers"
	"github.com/skrpld/NearBeee/internal/transport/grpc/interceptors"
	"github.com/skrpld/NearBeee/pkg/api/auth"
	"github.com/skrpld/NearBeee/pkg/api/posts"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type NearBeeeServerConfig struct {
	Host   string `env:"SERVER_HOST" env-default:"localhost" mapstructure:"SERVER_HOST"`
	Port   int    `env:"SERVER_PORT" env-default:"50050" mapstructure:"SERVER_PORT"`
	Secret string `env:"SECRET" env-default:"secret" mapstructure:"SECRET"`
}

type NearBeeeServer struct {
	cfg        NearBeeeServerConfig
	grpcServer *grpc.Server
	listener   net.Listener
	logger     logger.Logger
}

func NewNearBeeeServer(cfg NearBeeeServerConfig, ctx *context.Context, mongoDB *mongodb.MongoDB, postgresDB *postgres.PostgresDB, logger logger.Logger) (*NearBeeeServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	repo := repository.NewNearBeeeRepository(ctx, postgresDB, mongoDB)
	authSrv := service.NewAuthService(repo, cfg.Secret)
	authController := controllers.NewAuthController(authSrv)

	postsSrv := service.NewPostsService(repo)
	postsController := controllers.NewPostsController(postsSrv)

	authInterceptor := interceptors.NewAuthInterceptorHandler(authSrv).AuthInterceptor
	loggerInterceptor := interceptors.LoggerInterceptor

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(authInterceptor(), loggerInterceptor(logger)),
	}

	grpcServer := grpc.NewServer(opts...)
	auth.RegisterAuthServiceServer(grpcServer, authController)
	posts.RegisterPostsServiceServer(grpcServer, postsController)

	return &NearBeeeServer{cfg, grpcServer, lis, logger}, nil
}

func (s *NearBeeeServer) Start() error {
	s.logger.With(zap.Time("started_at", time.Now())).Info(fmt.Sprintf("Server started on %s:%d", s.cfg.Host, s.cfg.Port))
	return s.grpcServer.Serve(s.listener)
}

func (s *NearBeeeServer) Stop() {
	s.grpcServer.GracefulStop()
	s.logger.With(zap.Time("stopped_at", time.Now())).Info("Server stopped")
}
