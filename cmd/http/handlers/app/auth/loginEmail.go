package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/services/auth/loginEmail"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateLoginEmailHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.LoginEmail
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeLoginEmail(model)

		l := logger.NewLogBuilder()
		handler := loginEmail.New(loginEmail.NewModel(model.Email, model.Password), nil, l)

		return request.SendResponse[loginEmail.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
			c.SetCookie(request.EncryptAuthenticationCookie(model.(string)))

			return nil
		}, false)
	}
}
