package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/KarolosLykos/sportsnews/internal/article"
	httperrors "github.com/KarolosLykos/sportsnews/internal/utils/http_errors"
	"github.com/KarolosLykos/sportsnews/internal/utils/logger"
)

type articleHandler struct {
	logger logger.Logger
	uc     article.UseCase
}

func NewArticleHandler(logger logger.Logger, uc article.UseCase) *articleHandler {
	return &articleHandler{
		logger: logger,
		uc:     uc,
	}
}

func (h *articleHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		a, err := h.uc.GetByID(c.Request().Context(), id)
		if err != nil {
			return httperrors.ErrorResponse(c, err)
		}

		return c.JSON(http.StatusOK, a.ToRest())
	}
}

func (h *articleHandler) List() echo.HandlerFunc {
	return func(c echo.Context) error {
		articles, err := h.uc.List(c.Request().Context())
		if err != nil {
			return httperrors.ErrorResponse(c, err)
		}

		return c.JSON(http.StatusOK, articles.ToRest())
	}
}
