package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KarolosLykos/sportsnews/config"
	"github.com/KarolosLykos/sportsnews/domain"
	mock "github.com/KarolosLykos/sportsnews/internal/article/mock"
	httperrors "github.com/KarolosLykos/sportsnews/internal/utils/http_errors"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

func TestArticleHandler_GetByID(t *testing.T) {
	log := getLogger()

	a := &domain.Article{
		ID:        "6406083ea019b8815f689907",
		ArticleID: "123",
		TeamID:    "team",
	}

	tt := []struct {
		name string
		id   string
		stub func(uc *mock.MockUseCase)
		code int
		err  error
	}{
		{
			name: "internal server error",
			id:   "6406083ea019b8815f689907",
			stub: func(uc *mock.MockUseCase) {
				uc.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			code: http.StatusInternalServerError,
			err:  httperrors.ErrInternal,
		},
		{
			name: "not found",
			id:   "6406083ea019b8815f689907",
			stub: func(uc *mock.MockUseCase) {
				uc.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).
					Return(nil, errors.New("no documents in result"))
			},
			code: http.StatusNotFound,
			err:  httperrors.ErrNotFound,
		},
		{
			name: "ok",
			id:   "6406083ea019b8815f689907",
			stub: func(uc *mock.MockUseCase) {
				uc.EXPECT().GetByID(gomock.Any(), gomock.Any()).Times(1).
					Return(a, nil)
			},
			code: http.StatusOK,
			err:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUseCase(ctrl)

			tc.stub(uc)
			h := NewArticleHandler(log, uc)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/"+tc.id, nil)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			getByID := h.GetByID()
			err := getByID(c)

			if err != nil && tc.err != nil {
				assert.Equal(t, tc.code, rec.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					res := &domain.ArticleRest{}
					err := json.NewDecoder(rec.Body).Decode(res)
					require.NoError(t, err)

					assert.Equal(t, "success", res.Status)
					assert.Equal(t, res.Data, a)
				}
			}
		})
	}
}

func TestArticleHandler_List(t *testing.T) {
	log := getLogger()

	a := &domain.Articles{
		Total: 3,
		Articles: []*domain.Article{
			{
				ID:        "6406083ea019b8815f689907",
				ArticleID: "1",
				TeamID:    "team",
			},
			{
				ID:        "6406083ea019b8815f689907",
				ArticleID: "2",
				TeamID:    "team",
			},
			{
				ID:        "6406083ea019b8815f689907",
				ArticleID: "3",
				TeamID:    "team",
			},
		},
	}

	tt := []struct {
		name string
		stub func(uc *mock.MockUseCase)
		code int
		err  error
	}{
		{
			name: "internal server error",
			stub: func(uc *mock.MockUseCase) {
				uc.EXPECT().List(gomock.Any()).Times(1).
					Return(nil, errors.New("something went wrong"))
			},
			code: http.StatusInternalServerError,
			err:  httperrors.ErrInternal,
		},
		{
			name: "ok",
			stub: func(uc *mock.MockUseCase) {
				uc.EXPECT().List(gomock.Any()).Times(1).
					Return(a, nil)
			},
			code: http.StatusOK,
			err:  nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mock.NewMockUseCase(ctrl)

			tc.stub(uc)
			h := NewArticleHandler(log, uc)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/articles/", nil)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			list := h.List()
			err := list(c)

			if err != nil && tc.err != nil {
				assert.Equal(t, tc.code, rec.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.code, rec.Code)

				if rec.Code == http.StatusOK {
					res := &domain.ArticlesRest{}
					err := json.NewDecoder(rec.Body).Decode(res)
					require.NoError(t, err)

					assert.Equal(t, "success", res.Status)
					assert.Equal(t, a.Total, res.Metadata.Total)
					assert.Equal(t, a.Articles, res.Data)
				}
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
