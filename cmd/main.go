package main

import (
	"context"
	"log"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/internal/server"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
	"github.com/KarolosLykos/sportsnews/internal/utils/mongodb"
	"github.com/KarolosLykos/sportsnews/internal/utils/redis"
)

func main() {
	ctx := context.Background()

	// Parse configuration.
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalln(err)
	}

	// Initialise Logger.
	appLogger := logger.Init(cfg)
	appLogger.Infof(ctx, "Dev: %v, Level: %s", cfg.Dev, cfg.Logger.LogLevel)

	// Connect to mongo.
	mongoDBConn, err := mongodb.New(ctx, cfg)
	if err != nil {
		appLogger.Fatal(ctx, err, "could not connect to mongodb")
	}

	defer mongoDBConn.Disconnect(ctx) //nolint:errCheck // TODO: check if all connections have been closed.

	appLogger.Infof(ctx, "connected to mongoDB: %v", mongoDBConn.NumberSessionsInProgress())

	// Connect to redis.
	redisClient := redis.New(cfg)
	defer redisClient.Close()

	appLogger.Infof(ctx, "connected to redis")

	// Init server.
	s := server.New(cfg, appLogger, mongoDBConn, redisClient)

	// Run server.
	if err = s.Run(); err != nil {
		appLogger.Fatal(ctx, err)
	}
}
