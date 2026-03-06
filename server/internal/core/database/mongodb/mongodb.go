package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDBConfig struct {
	DBName string `env:"MONGODB_DB" env-default:"nearbeee" mapstructure:"MONGODB_DB"`
	Host   string `env:"MONGODB_HOST" env-default:"localhost" mapstructure:"MONGODB_HOST"`
	Port   string `env:"MONGODB_PORT" env-default:"27017" mapstructure:"MONGODB_PORT"`
}
type MongoDB struct {
	*mongo.Database
}

func NewMongoDB(cfg MongoDBConfig, ctx context.Context) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://" + cfg.Host + ":" + cfg.Port))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.DBName)

	return &MongoDB{db}, nil
}

func (m *MongoDB) Close() error {
	if m.Client() == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := m.Client().Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
