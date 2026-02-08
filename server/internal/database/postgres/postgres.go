package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresConfig struct {
	DBName   string `env:"POSTGRES_DB" env-default:"blog"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	User     string `env:"POSTGRES_USER" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable"`
}

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB(cfg PostgresConfig, ctx context.Context) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	return &PostgresDB{db}, nil
}
