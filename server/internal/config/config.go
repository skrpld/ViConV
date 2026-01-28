package config

import (
	"viconv/internal/database/mongodb"
	"viconv/internal/database/postgres"
	"viconv/internal/transport/grpc/servers"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	servers.ViconvServerConfig
	mongodb.MongoDBConfig
	postgres.PostgresConfig
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
