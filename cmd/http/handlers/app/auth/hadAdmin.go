package auth

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/services/auth/adminExists"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AdminExistsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := logger.NewLogBuilder()
		handler := adminExists.New(l)

		return request.SendResponse[interface{}](handler, c, http.StatusOK, l, nil, false)
	}
}
