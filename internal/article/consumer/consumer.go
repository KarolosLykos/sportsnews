package consumer

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/domain"
	"github.com/KarolosLykos/sportsnews/internal/article"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

var (
	ErrGetByID   = errors.New("consumer: getByID")
	ErrList      = errors.New("consumer: list")
	ErrBadStatus = errors.New("bad status code")
)

type HullCityConsumer struct {
	cfg        *config.Config
	logger     logger.Logger
	client     *http.Client
	repository article.Repository
	cache      article.Cache
}

func NewHullCityConsumer(
	cfg *config.Config,
	logger logger.Logger,
	client *http.Client,
	repository article.Repository,
	cache article.Cache,
) *HullCityConsumer {
	return &HullCityConsumer{
		cfg:        cfg,
		logger:     logger,
		client:     client,
		repository: repository,
		cache:      cache,
	}
}

func (c *HullCityConsumer) GetByID(ctx context.Context, id string) (*domain.HullArticleInformation, error) {
	uri := c.cfg.Consumer.HullConsumer.SingleURL + "?id=" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGetByID, err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGetByID, err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w:%v", ErrBadStatus, res.Status)
	}

	defer res.Body.Close()

	hullArticle := &domain.HullArticleInformation{}
	if err = xml.NewDecoder(res.Body).Decode(hullArticle); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrGetByID, err)
	}

	return hullArticle, nil
}
func (c *HullCityConsumer) List(ctx context.Context) (*domain.HullArticles, error) {
	uri := c.cfg.Consumer.HullConsumer.ListURL + "?count=" + strconv.Itoa(c.cfg.Consumer.HullConsumer.Count)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w:%v", ErrBadStatus, res.Status)
	}

	defer res.Body.Close()

	hullArticles := &domain.HullArticles{}
	if err = xml.NewDecoder(res.Body).Decode(hullArticles); err != nil {
		return nil, fmt.Errorf("%w:%v", ErrList, err)
	}

	return hullArticles, nil
}

func (c *HullCityConsumer) Consume(ctx context.Context) {
	hullArticles, err := c.List(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return
	}

	for _, h := range hullArticles.NewsletterNewsItems.NewsletterNewsItem {
		item, err := c.GetByID(ctx, h.NewsArticleID)
		if err != nil {
			c.logger.Warn(ctx, err)
			continue
		}

		a := h.ToDomain(
			hullArticles.ClubName,
			hullArticles.ClubWebsiteURL,
			item.NewsArticle.BodyText,
			item.NewsArticle.Subtitle,
		)

		updatedArticle, err := c.repository.Upsert(ctx, a)
		if err != nil {
			c.logger.Warn(ctx, err)
			continue
		}

		if err = c.cache.Set(ctx, updatedArticle); err != nil {
			c.logger.Warn(ctx, err)
		}
	}
}
