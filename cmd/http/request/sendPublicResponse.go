package request

import (
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SendPublicResponse[T any, F any, K any](handler pkg.Job[T, F, K], context echo.Context, status int, callback func(c echo.Context, model interface{}) error, gracefulFail bool) error {
	model, err := handler.Handle()

	if err != nil {
		appError, ok := err.(publicApiError.PublicApiError)
		if ok {
			s := http.StatusInternalServerError
			if appError.Status() == publicApiError.ValidationError {
				s = http.StatusUnprocessableEntity
			} else if appError.Status() == publicApiError.NotFoundError {
				s = http.StatusNotFound
			} else if appError.Status() == publicApiError.ApplicationError {
				s = http.StatusBadRequest
			}

			return context.JSON(s, appError.Data())
		} else {
			return context.JSON(http.StatusInternalServerError, map[string]string{
				"data": err.Error(),
			})
		}
	}

	if err := callCallback(context, model, callback); err != nil {
		if err := context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		}); err != nil {
			return err
		}

		return nil
	}

	if err := context.JSON(status, model); err != nil {
		return err
	}

	return nil
}
