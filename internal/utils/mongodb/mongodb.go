package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/KarolosLykos/sportsnews/config"
)

var ErrMongoConn = errors.New("mongo connection error")

// New creates a new connection to mongoDb.
func New(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.MongoDB.Host, cfg.MongoDB.Port)).
			SetAuth(options.Credential{
				Username: cfg.MongoDB.Username,
				Password: cfg.MongoDB.Password,
			}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrMongoConn, err)
	}

	if err = client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrMongoConn, err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, err
}
