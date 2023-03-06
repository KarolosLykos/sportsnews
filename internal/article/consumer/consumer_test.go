package consumer

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/domain"
	mock "github.com/KarolosLykos/sportsnews/internal/article/mock"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

func TestHullCityConsumer_GetByID(t *testing.T) {
	log := getLogger()

	cfg := &config.Config{Consumer: config.ConsumerConfig{HullConsumer: config.HullConsumer{
		SingleURL: "single",
	}}}

	testArticle := &domain.HullArticleInformation{
		NewsArticle: domain.HullArticle{
			ArticleURL: "test.com",
			Title:      "test title",
			Subtitle:   "test subtitle",
			BodyText:   "test content",
		},
	}

	tt := []struct {
		name      string
		id        string
		responder httpmock.Responder
		err       error
	}{
		{
			name:      "ok",
			id:        "6405f896a019b8815f6892c9",
			responder: httpmock.NewStringResponder(http.StatusOK, testXMLSingle),
			err:       nil,
		},
		{
			name:      "invalid xml response",
			id:        "6405f896a019b8815f6892c9",
			responder: httpmock.NewStringResponder(http.StatusOK, "{}"),
			err:       ErrGetByID,
		},
		{
			name:      "bad status code",
			id:        "6405f896a019b8815f6892c9",
			responder: httpmock.NewStringResponder(http.StatusBadRequest, "{}"),
			err:       ErrBadStatus,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)
			cache := mock.NewMockCache(ctrl)

			httpmock.RegisterResponder(http.MethodGet, "single?id="+tc.id, tc.responder)

			c := NewHullCityConsumer(cfg, log, &http.Client{}, repo, cache)

			a, err := c.GetByID(context.Background(), tc.id)
			if err != nil && tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, a)
				assert.Equal(t, testArticle, a)
			}
		})
	}
}

func TestHullCityConsumer_List(t *testing.T) {
	log := getLogger()

	cfg := &config.Config{Consumer: config.ConsumerConfig{HullConsumer: config.HullConsumer{
		ListURL: "list",
		Count:   3,
	}}}

	testArticles := &domain.HullArticles{
		ClubName:       "clubname",
		ClubWebsiteURL: "club.com",
		NewsletterNewsItems: struct {
			NewsletterNewsItem []domain.HullArticle `xml:"NewsletterNewsItem"`
		}{
			[]domain.HullArticle{
				{ArticleURL: "test1.com", Title: "test title1"},
				{ArticleURL: "test2.com", Title: "test title2"},
				{ArticleURL: "test3.com", Title: "test title3"},
			},
		},
	}

	tt := []struct {
		name      string
		responder httpmock.Responder
		err       error
	}{
		{
			name:      "ok",
			responder: httpmock.NewStringResponder(http.StatusOK, testXMLList),
			err:       nil,
		},
		{
			name:      "invalid xml response",
			responder: httpmock.NewStringResponder(http.StatusOK, "{}"),
			err:       ErrList,
		},
		{
			name:      "bad status code",
			responder: httpmock.NewStringResponder(http.StatusBadRequest, "{}"),
			err:       ErrBadStatus,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)
			cache := mock.NewMockCache(ctrl)

			httpmock.RegisterResponder(http.MethodGet, "list?count=3", tc.responder)

			c := NewHullCityConsumer(cfg, log, &http.Client{}, repo, cache)

			a, err := c.List(context.Background())
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

var (
	testXMLSingle = `<NewsArticleInformation>
<ClubName>Hull City</ClubName>
<ClubWebsiteURL>test.com</ClubWebsiteURL>
<NewsArticle>
<ArticleURL>test.com</ArticleURL>
<Title>test title</Title>
<BodyText>test content</BodyText>
<Subtitle>test subtitle</Subtitle>
</NewsArticle>
</NewsArticleInformation>`

	testXMLList = `<NewListInformation>
<ClubName>clubname</ClubName>
<ClubWebsiteURL>club.com</ClubWebsiteURL>
<NewsletterNewsItems>
<NewsletterNewsItem>
<ArticleURL>test1.com</ArticleURL>
<Title>test title1</Title>
</NewsletterNewsItem>
<NewsletterNewsItem>
<ArticleURL>test2.com</ArticleURL>
<Title>test title2</Title>
</NewsletterNewsItem>
<NewsletterNewsItem>
<ArticleURL>test3.com</ArticleURL>
<Title>test title3</Title>
</NewsletterNewsItem>
</NewsletterNewsItems>
</NewListInformation>`
)
