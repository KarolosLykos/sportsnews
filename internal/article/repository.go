package article

import (
	"context"

	"github.com/KarolosLykos/sportsnews/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	List(ctx context.Context) (*domain.Articles, error)
	Upsert(ctx context.Context, article *domain.Article) (*domain.Article, error)
}

type Cache interface {
	Get(ctx context.Context, id string) (*domain.Article, error)
	Set(ctx context.Context, article *domain.Article) error
}
