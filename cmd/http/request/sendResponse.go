package request

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type DevErrorResponse struct {
	StackTrace string `json:"stackTrace"`
	Error      string `json:"error"`
}

type ErrorResponse[T any] struct {
	Data T `json:"data"`
}

func SendResponse[T any, F any, K any](handler pkg.Job[T, F, K], context echo.Context, status int, logger logger.LogBuilder, callback func(c echo.Context, model interface{}) error) error {
	model, err := handler.Handle()

	if err != nil {
		validationError, ok := err.(appErrors.AppError[map[string]string])
		fmt.Println(err, validationError, ok)
		if ok {
			if err := flushLogger(logger, "info", context); err != nil {
				return err
			}
			return context.JSON(http.StatusUnprocessableEntity, ErrorResponse[map[string]string]{
				Data: validationError.Data(),
			})
		}

		otherError, ok := err.(appErrors.AppError[struct{}])
		if ok {
			if otherError.Type() == appErrors.AUTHENTICATION_ERROR {
				if err := flushLogger(logger, "info", context); err != nil {
					return err
				}
				return context.JSON(http.StatusForbidden, ErrorResponse[map[string]string]{
					Data: map[string]string{
						"unauthenticated": "You are not authenticated",
					},
				})
			} else if otherError.Type() == appErrors.AUTHORIZATION_ERROR {
				if err := flushLogger(logger, "info", context); err != nil {
					return err
				}
				return context.JSON(http.StatusUnauthorized, ErrorResponse[map[string]string]{
					Data: map[string]string{
						"unauthorized": "You are not authorized",
					},
				})
			} else if otherError.Type() == appErrors.NOT_FOUND_ERROR {
				if err := flushLogger(logger, "info", context); err != nil {
					return err
				}
				return context.JSON(http.StatusNotFound, ErrorResponse[map[string]string]{
					Data: map[string]string{
						"notExists": "The requested resource does not exist.",
					},
				})
			} else if otherError.Type() == appErrors.USER_UNCOFIRMED {
				if err := flushLogger(logger, "info", context); err != nil {
					return err
				}
				return context.JSON(http.StatusNotFound, ErrorResponse[map[string]string]{
					Data: map[string]string{
						"userUnconfirmed": "The user is unconfirmed",
					},
				})
			}

			if os.Getenv("APP_ENV") != "prod" {
				if err := flushLogger(logger, "error", context); err != nil {
					return err
				}
				return context.JSON(http.StatusInternalServerError, ErrorResponse[DevErrorResponse]{
					Data: DevErrorResponse{
						StackTrace: otherError.StackTrace(),
						Error:      otherError.Error(),
					},
				})
			}
		}

		if err := flushLogger(logger, "error", context); err != nil {
			return err
		}
		return context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		})
	}

	if err := callCallback(context, model, callback); err != nil {
		if err := flushLogger(logger, "info", context); err != nil {
			return err
		}

		return context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		})
	}

	if err := flushLogger(logger, "info", context); err != nil {
		return err
	}

	return context.JSON(status, model)
}

func flushLogger(logger logger.LogBuilder, t string, context echo.Context) error {
	if err := logger.Flush(t); err != nil {
		return context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		})
	}

	return nil
}

func callCallback(c echo.Context, model interface{}, cb func(c echo.Context, model interface{}) error) error {
	if cb != nil {
		return cb(c, model)
	}

	return nil
}
