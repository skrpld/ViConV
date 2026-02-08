package servers

import (
	"context"
	"fmt"
	"net"
	"time"
	"viconv/internal/database/mongodb"
	"viconv/internal/database/postgres"
	"viconv/internal/logger"
	"viconv/internal/repository"
	"viconv/internal/service"
	"viconv/internal/transport/grpc/controllers"
	"viconv/internal/transport/grpc/interceptors"
	"viconv/pkg/api/auth"
	"viconv/pkg/api/posts"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ViconvServerConfig struct {
	Host   string `env:"SERVER_HOST" env-default:"localhost" mapstructure:"SERVER_HOST"`
	Port   int    `env:"SERVER_PORT" env-default:"50050" mapstructure:"SERVER_PORT"`
	Secret string `env:"SECRET" env-default:"secret" mapstructure:"SECRET"`
}

type ViconvServer struct {
	cfg        ViconvServerConfig
	grpcServer *grpc.Server
	listener   net.Listener
	logger     logger.Logger
}

func NewViconvServer(cfg ViconvServerConfig, ctx *context.Context, mongoDB *mongodb.MongoDB, postgresDB *postgres.PostgresDB, logger logger.Logger) (*ViconvServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	repo := repository.NewViconvRepository(ctx, postgresDB, mongoDB)
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

	return &ViconvServer{cfg, grpcServer, lis, logger}, nil
}

func (s *ViconvServer) Start() error {
	s.logger.With(zap.Time("started_at", time.Now())).Info(fmt.Sprintf("Server started on %s:%d", s.cfg.Host, s.cfg.Port))
	return s.grpcServer.Serve(s.listener)
}

func (s *ViconvServer) Stop() {
	s.grpcServer.GracefulStop()
	s.logger.With(zap.Time("stopped_at", time.Now())).Info("Server stopped")
}
