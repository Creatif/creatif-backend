package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/services/auth/loginApi"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateLoginApiHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.LoginApi
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeLoginApi(model)

		apiKey := c.Request().Header.Get("X-CREATIF-API-KEY")
		projectId := c.Request().Header.Get("X-CREATIF-PROJECT-ID")

		l := logger.NewLogBuilder()
		handler := loginApi.New(loginApi.NewModel(model.Email, model.Password, apiKey, projectId, model.Session), nil, l)

		return request.SendResponse[loginApi.Model](handler, c, http.StatusOK, l, nil)
	}
}
