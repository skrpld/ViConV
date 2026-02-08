package servers

import (
	"context"
	"fmt"
	"net"
	"viconv/internal/database/mongodb"
	"viconv/internal/logger"
	"viconv/internal/repository"
	"viconv/internal/service"
	"viconv/internal/transport/grpc/controllers"
	"viconv/internal/transport/grpc/interceptors"
	"viconv/pkg/api/auth"

	"google.golang.org/grpc"
)

type ViconvServerConfig struct {
	Host   string `env:"HOST" env-default:"localhost"`
	Port   int    `env:"PORT" env-default:"50050"`
	Secret string `env:"SECRET" env-default:"secret"`
}

type ViconvServer struct {
	cfg        ViconvServerConfig
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewViconvServer(cfg ViconvServerConfig, ctx *context.Context, db *mongodb.DB, logger logger.Logger) (*ViconvServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	repo := repository.NewViconvRepository(ctx, db)
	authSrv := service.NewAuthService(repo, cfg.Secret)
	authController := controllers.NewAuthController(authSrv)

	authInterceptor := interceptors.NewAuthInterceptorHandler(authSrv).AuthInterceptor
	loggerInterceptor := interceptors.LoggerInterceptor

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(authInterceptor(), loggerInterceptor(logger)),
	}

	grpcServer := grpc.NewServer(opts...)
	auth.RegisterAuthServiceServer(grpcServer, authController)

	return &ViconvServer{cfg, grpcServer, lis}, nil
}

func (s *ViconvServer) Start() error {
	fmt.Printf("Server started on %s:%d\n", s.cfg.Host, s.cfg.Port)
	return s.grpcServer.Serve(s.listener)
}

func (s *ViconvServer) Stop() {
	s.grpcServer.GracefulStop()
}
