package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/services/auth/loginApi"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoginHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.LoginApi
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeLoginApi(model)

		handler := loginApi.New(loginApi.NewModel(model.Email, model.Password), nil)

		return request.SendResponse[loginApi.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
			c.SetCookie(request.EncryptApiAuthenticationCookie(model.(string)))

			return nil
		}, false)
	}
}
