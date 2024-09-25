package auth

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/services/auth/adminExists"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AdminExistsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		handler := adminExists.New()

		return request.SendResponse[interface{}](handler, c, http.StatusOK, nil, false)
	}
}
