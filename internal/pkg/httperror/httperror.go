package httperror

import (
	"net/http"

	"github.com/ell1jah/bmstu_web/model"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func HandleError(err error) *echo.HTTPError {
	causeErr := errors.Cause(err)
	switch {
	case errors.Is(causeErr, model.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, model.ErrNotFound.Error())
	case errors.Is(causeErr, model.ErrBadRequest):
		return echo.NewHTTPError(http.StatusBadRequest, model.ErrBadRequest.Error())
	case errors.Is(causeErr, model.ErrPermissionDenied):
		return echo.NewHTTPError(http.StatusForbidden, model.ErrPermissionDenied.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, causeErr.Error())
	}
}
