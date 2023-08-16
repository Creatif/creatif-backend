package request

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
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

func SendResponse[T any, F any, K any](handler pkg.Job[T, F, K], context echo.Context, status int) error {
	createdModel, err := handler.Handle()

	if err != nil {
		validationError, ok := err.(appErrors.AppError[map[string]string])
		if ok {
			return context.JSON(http.StatusUnprocessableEntity, ErrorResponse[map[string]string]{
				Data: validationError.Data(),
			})
		}

		otherError, ok := err.(appErrors.AppError[struct{}])
		if ok {
			if otherError.Type() == appErrors.AUTHENTICATION_ERROR {
				return context.JSON(http.StatusForbidden, ErrorResponse[string]{
					Data: "Unauthenticated",
				})
			} else if otherError.Type() == appErrors.AUTHORIZATION_ERROR {
				return context.JSON(http.StatusUnauthorized, ErrorResponse[string]{
					Data: "Unauthorized",
				})
			}

			if os.Getenv("APP_ENV") != "prod" {
				return context.JSON(http.StatusInternalServerError, ErrorResponse[DevErrorResponse]{
					Data: DevErrorResponse{
						StackTrace: otherError.StackTrace(),
						Error:      otherError.Error(),
					},
				})
			}
		}

		return context.JSON(http.StatusInternalServerError, ErrorResponse[string]{
			Data: "Internal server error",
		})
	}

	return context.JSON(status, createdModel)
}
