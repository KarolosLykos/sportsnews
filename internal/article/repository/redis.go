package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/domain"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

var (
	ErrMarshal   = errors.New("cache: couldn't marshal")
	ErrUnMarshal = errors.New("cache: couldn't unmarshal")
	ErrSet       = errors.New("cache: couldn't set value")
	ErrGet       = errors.New("cache: couldn't get value")
)

type Cache struct {
	cfg    *config.Config
	logger logger.Logger
	client *redis.Client
}

func NewCacheRepository(cfg *config.Config, logger logger.Logger, client *redis.Client) *Cache {
	return &Cache{cfg: cfg, logger: logger, client: client}
}

func (c Cache) Get(ctx context.Context, id string) (*domain.Article, error) {
	res, err := c.client.Get(ctx, getKey(c.cfg.Redis.KeyPrefix, id)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGet, err)
	}

	article := &domain.Article{}
	if err = json.Unmarshal(res, article); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrUnMarshal, err)
	}

	return article, nil
}

func (c Cache) Set(ctx context.Context, article *domain.Article) error {
	articleB, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("%w:%v", ErrMarshal, err)
	}

	if err = c.client.Set(ctx, getKey(c.cfg.Redis.KeyPrefix, article.ID), articleB, c.cfg.Redis.Expiration).
		Err(); err != nil {
		return fmt.Errorf("%w:%v", ErrSet, err)
	}

	return nil
}

func getKey(prefix, id string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}
