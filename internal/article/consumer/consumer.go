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

	maxWorkers = 30
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

// Consume consumes feeds from the Hull City Fc external provider.
func (c *HullCityConsumer) Consume(ctx context.Context) {
	hullArticles, err := c.List(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return
	}

	jobs := make(chan domain.HullArticle, len(hullArticles.NewsletterNewsItems.NewsletterNewsItem))

	if len(hullArticles.NewsletterNewsItems.NewsletterNewsItem) < maxWorkers {
		maxWorkers = len(hullArticles.NewsletterNewsItems.NewsletterNewsItem)
	}

	// Start workers to process each article item.
	for w := 0; w < maxWorkers; w++ {
		go c.worker(ctx, jobs, hullArticles.ClubName, hullArticles.ClubWebsiteURL)
	}

	// Send each article item to the job channel.
	for _, h := range hullArticles.NewsletterNewsItems.NewsletterNewsItem {
		jobs <- h
	}

	close(jobs)
}

func (c *HullCityConsumer) worker(ctx context.Context, jobs <-chan domain.HullArticle, clubName, clubURL string) {
	for j := range jobs {
		c.logger.Debugf(ctx, "processing job for article %s", j.NewsArticleID)

		item, err := c.GetByID(ctx, j.NewsArticleID)
		if err != nil {
			c.logger.Warn(ctx, err)
			continue
		}

		a := j.ToDomain(
			clubName,
			clubURL,
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
