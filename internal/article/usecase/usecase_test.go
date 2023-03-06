package usecase

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/domain"
	mock "github.com/KarolosLykos/sportsnews/internal/article/mock"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

func TestArticleUseCase_GetByID(t *testing.T) {
	log := getLogger()

	testArticle := &domain.Article{
		ID:          "6405f896a019b8815f6892c7",
		TeamID:      "Hull City",
		ClubURL:     "https://www.wearehullcity.co.uk",
		OptaMatchID: "g2322054",
		Title:       "Hall: ‘Really happy with our team performance’",
		Type:        []string{"Academy"},
		Teaser:      "Midfielder Sincere Hall was delighted with the team performance as the Under-21s defeated Sheffield Wednesday 1-0 at the MKM Stadium.",
		URL:         "https://www.wearehullcity.co.uk/news/2023/march/hall-really-happy-with-our-team-performance/",
		ImageURL:    "https://www.wearehullcity.co.uk/api/image/feedassets/f5582976-c069-4b12-9da4-e394c428deb3/Medium/sincere-hall.jpg",
		IsPublished: true,
	}

	tt := []struct {
		name      string
		id        string
		repoStub  func(repo *mock.MockRepository)
		cacheStub func(repo *mock.MockCache)
		err       error
	}{
		{
			name: "ok not cached",
			id:   "6405f896a019b8815f6892c9",
			repoStub: func(repo *mock.MockRepository) {
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(testArticle, nil)
			},
			cacheStub: func(cache *mock.MockCache) {
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
			err: nil,
		},
		{
			name: "ok cached",
			id:   "6405f896a019b8815f6892c9",
			repoStub: func(repo *mock.MockRepository) {
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(0)
			},
			cacheStub: func(cache *mock.MockCache) {
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(testArticle, nil)
				cache.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)
			},
			err: nil,
		},
		{
			name: "not found",
			id:   "asdasd",
			repoStub: func(repo *mock.MockRepository) {
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("no documents in result"))
			},
			cacheStub: func(cache *mock.MockCache) {
				cache.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(nil, redis.Nil)
				cache.EXPECT().Set(gomock.Any(), gomock.Any()).Times(0)
			},
			err: ErrGetByID,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)
			cache := mock.NewMockCache(ctrl)

			tc.repoStub(repo)
			tc.cacheStub(cache)

			uc := New(log, repo, cache)

			a, err := uc.GetByID(context.Background(), tc.id)
			if err != nil && tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, a)
				assert.Equal(t, a, testArticle)
			}
		})
	}
}

func TestArticleUseCase_List(t *testing.T) {
	log := getLogger()

	testArticles := &domain.Articles{
		Total: 1,
		Articles: []*domain.Article{
			{
				ID:          "6405f896a019b8815f6892c7",
				TeamID:      "Hull City",
				ClubURL:     "https://www.wearehullcity.co.uk",
				OptaMatchID: "g2322054",
				Title:       "Hall: ‘Really happy with our team performance’",
				Type:        []string{"Academy"},
				Teaser:      "Midfielder Sincere Hall was delighted with the team performance as the Under-21s defeated Sheffield Wednesday 1-0 at the MKM Stadium.",
				URL:         "https://www.wearehullcity.co.uk/news/2023/march/hall-really-happy-with-our-team-performance/",
				ImageURL:    "https://www.wearehullcity.co.uk/api/image/feedassets/f5582976-c069-4b12-9da4-e394c428deb3/Medium/sincere-hall.jpg",
				IsPublished: true,
			},
		},
	}

	tt := []struct {
		name     string
		repoStub func(repo *mock.MockRepository)
		err      error
	}{
		{
			name: "ok",
			repoStub: func(repo *mock.MockRepository) {
				repo.EXPECT().List(gomock.Any()).Times(1).Return(testArticles, nil)
			},
			err: nil,
		},
		{
			name: "generic err",
			repoStub: func(repo *mock.MockRepository) {
				repo.EXPECT().List(gomock.Any()).Times(1).Return(nil, errors.New("generic error"))
			},
			err: ErrList,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)
			cache := mock.NewMockCache(ctrl)

			tc.repoStub(repo)

			uc := New(log, repo, cache)

			a, err := uc.List(context.Background())
			if err != nil && tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, a)
				assert.Equal(t, testArticles, a)
			}
		})
	}
}

func getLogger() logger.Logger {
	cfg := &config.Config{}

	l := &logrus.Logger{
		Out:          io.Discard,
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
		Level:        logrus.InfoLevel,
		Formatter:    &logrus.JSONFormatter{},
	}

	return logger.New(cfg, l)
}
