package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDBConfig struct {
	DBName string `env:"MONGODB_DB" env-default:"viconv"`
	Host   string `env:"MONGODB_HOST" env-default:"localhost"`
	Port   string `env:"MONGODB_PORT" env-default:"27017"`
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
