package request

import (
	"creatif/pkg/app/services/publicApi/publicApiError"
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SendPublicResponse[T any, F any, K any](handler pkg.Job[T, F, K], context echo.Context, status int, lg logger.LogBuilder, callback func(c echo.Context, model interface{}) error, gracefulFail bool) error {
	model, err := handler.Handle()

	if err != nil {
		appError, ok := err.(publicApiError.PublicApiError)
		if ok {
			fmt.Println(appError)
			s := http.StatusInternalServerError
			if appError.Status() == publicApiError.ValidationError {
				s = http.StatusUnprocessableEntity
			} else if appError.Status() == publicApiError.NotFoundError {
				s = http.StatusNotFound
			} else if appError.Status() == publicApiError.ApplicationError {
				s = http.StatusBadRequest
			}

			return context.JSON(s, appError.Data())
		}
	}

	if err := callCallback(context, model, callback); err != nil {
		lg.Add("Callback error", err.Error())
		if err := flushLogger(lg, "info", context); err != nil {
			fmt.Println("Flush error: ", err)
			return err
		}

		if err := context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		}); err != nil {
			logger.Error(err.Error())
			return err
		}

		return nil
	}

	if err := flushLogger(lg, "info", context); err != nil {
		return err
	}

	if err := context.JSON(status, model); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}
