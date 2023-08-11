package request

import (
	pkg "creatif/pkg/lib"
	"creatif/pkg/lib/appErrors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	StackTrace string `json:"stackTrace"`
	Error      string `json:"error"`
}

func SendResponse[T any, F any, K any](handler pkg.Job[T, F, K], context echo.Context, status int) error {
	createdModel, err := handler.Handle()

	if err != nil {
		validationError, ok := err.(appErrors.AppError[map[string]string])
		if ok {
			return context.JSON(http.StatusUnprocessableEntity, validationError.Data())
		}

		otherError, ok := err.(appErrors.AppError[struct{}])
		if ok {
			if otherError.Type() == appErrors.AUTHENTICATION_ERROR {
				return context.JSON(http.StatusForbidden, "Forbidden")
			} else if otherError.Type() == appErrors.AUTHORIZATION_ERROR {
				return context.JSON(http.StatusUnauthorized, "Unauthorized")
			}

			return context.JSON(http.StatusInternalServerError, Response{
				StackTrace: otherError.StackTrace(),
				Error:      otherError.Error(),
			})
		}

		return context.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return context.JSON(status, createdModel)
}
