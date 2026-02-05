package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	DBName   string `env:"POSTGRES_DB" env-default:"viconv" mapstructure:"POSTGRES_DB"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432" mapstructure:"POSTGRES_PORT"`
	Host     string `env:"POSTGRES_HOST" default:"localhost" mapstructure:"POSTGRES_HOST"`
	User     string `env:"POSTGRES_USER" env-default:"postgres" mapstructure:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123" mapstructure:"POSTGRES_PASSWORD"`
	SSLMode  string `env:"POSTGRES_SSLMODE" env-default:"disable" mapstructure:"POSTGRES_SSLMODE"`
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

func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
