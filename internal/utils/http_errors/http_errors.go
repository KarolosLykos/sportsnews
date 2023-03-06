package httperrors

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	ErrInternal   = errors.New("something went wrong")
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)

type RestErr struct {
	Status int
	Err    string      `json:"error"`
	Cause  interface{} `json:"cause"`
}

// NewRestError returns new rest error.
func NewRestError(status int, err string, causes interface{}) *RestErr {
	return &RestErr{
		Status: status,
		Err:    err,
		Cause:  causes,
	}
}

func ErrorResponse(ctx echo.Context, err error) error {
	restErr := parseError(err)
	return ctx.JSON(restErr.Status, restErr)
}

func parseError(err error) *RestErr {
	switch {
	case strings.Contains(err.Error(), "no documents in result"):
		return NewRestError(http.StatusNotFound, ErrNotFound.Error(), err.Error())
	case strings.Contains(err.Error(), "provided hex string is not a valid ObjectID"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err.Error())
	}

	return NewRestError(http.StatusInternalServerError, ErrInternal.Error(), err.Error())
}
