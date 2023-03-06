package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/KarolosLykos/sportsnews/domain"
	"github.com/KarolosLykos/sportsnews/internal/article"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

var (
	ErrGetByID = errors.New("usecase: getByID")
	ErrList    = errors.New("usecase: list")
)

type articleUseCase struct {
	logger     logger.Logger
	repository article.Repository
	cache      article.Cache
}

func New(logger logger.Logger, repository article.Repository, cache article.Cache) *articleUseCase {
	return &articleUseCase{
		logger:     logger,
		repository: repository,
		cache:      cache,
	}
}

func (u *articleUseCase) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	cached, err := u.cache.Get(ctx, id)
	if err != nil {
		u.logger.Warnf(ctx, err, "could not get cached article with id: %s", id)
	}

	if cached != nil {
		return cached, nil
	}

	art, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGetByID, err)
	}

	if err = u.cache.Set(ctx, art); err != nil {
		u.logger.Warnf(ctx, err, "could not set article with id: %s", id)
	}

	return art, nil
}

func (u *articleUseCase) List(ctx context.Context) (*domain.Articles, error) {
	articles, err := u.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	return articles, nil
}
