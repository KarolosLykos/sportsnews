package article

import (
	"context"

	"github.com/KarolosLykos/sportsnews/domain"
)

type UseCase interface {
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	List(ctx context.Context) (*domain.Articles, error)
}
