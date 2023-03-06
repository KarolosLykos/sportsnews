package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/internal/article"
	"github.com/KarolosLykos/sportsnews/internal/article/consumer"
	v1 "github.com/KarolosLykos/sportsnews/internal/article/delivery/http/v1"
	"github.com/KarolosLykos/sportsnews/internal/article/repository"
	"github.com/KarolosLykos/sportsnews/internal/article/usecase"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

type Server struct {
	cfg    *config.Config
	logger logger.Logger

	httpServer *echo.Echo

	mongoDB     *mongo.Client
	redisClient *redis.Client
}

func New(
	cfg *config.Config,
	logger logger.Logger,
	mongoDB *mongo.Client,
	redisClient *redis.Client,
) *Server {
	return &Server{
		cfg:         cfg,
		logger:      logger,
		mongoDB:     mongoDB,
		redisClient: redisClient,
	}
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create new mongo repository
	mongoRepo := repository.NewMongoRepository(s.mongoDB, s.logger)
	// Create new redis cache.
	redisCache := repository.NewCacheRepository(s.cfg, s.logger, s.redisClient)
	// Create new article useCase.
	articleUC := usecase.New(s.logger, mongoRepo, redisCache)
	// Create new hullCity consumer.
	hullCityConsumer := consumer.NewHullCityConsumer(s.cfg, s.logger, http.DefaultClient, mongoRepo, redisCache)

	// Setup cron.
	cron := gocron.NewScheduler(time.UTC)
	j, err := cron.Every(s.cfg.Consumer.HullConsumer.Frequency).Do(hullCityConsumer.List, ctx)
	if err != nil {
		fmt.Printf("Job: %v, Error: %v", j, err)
	}
	cron.StartAsync()

	s.httpServer = s.createHTTP(articleUC)
	go func() {
		s.logger.Infof(ctx, "http server listening on port: %s", s.cfg.HTTP.Port)
		if err := s.httpServer.Start(s.cfg.HTTP.Port); err != nil {
			s.logger.Warn(ctx, err, "http server error ")
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case q := <-quit:
		s.logger.Infof(ctx, "received signal: %v", q)
	case <-ctx.Done():
		s.logger.Error(ctx, ctx.Err(), "context done")
	}

	// Gracefully shutdown servers.

	s.gracefullyShutdown(ctx)
	cron.Stop()

	return nil
}

// createHTTP creates new instance of Echo.
func (s *Server) createHTTP(
	uc article.UseCase,
) *echo.Echo {
	e := echo.New()
	e.Logger.SetOutput(io.Discard)
	e.Use(middleware.Recover())

	articleHandler := v1.NewArticleHandler(s.logger, uc)

	group := e.Group("/api/v1/articles")
	group.GET("/:id", articleHandler.GetByID())
	group.GET("", articleHandler.List())

	return e
}

// gracefullyShutdown gracefully shutdown servers.
func (s *Server) gracefullyShutdown(ctx context.Context) {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Warn(ctx, err, "http shutdown server")
	}

	s.logger.Info(ctx, "service exited gracefully")
}
