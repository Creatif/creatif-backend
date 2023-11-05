package auth

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/auth/loginEmail"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func CreateLoginEmailHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.LoginEmail
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeLoginEmail(model)

		l := logger.NewLogBuilder()
		handler := loginEmail.New(loginEmail.NewModel(model.Email, model.Password), auth.NewFrontendAuthentication(), l)

		return request.SendResponse[loginEmail.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) {
			encryptedUser := model.([]byte)

			cookie := new(http.Cookie)
			cookie.Name = "authentication"
			cookie.HttpOnly = true
			cookie.Secure = true
			cookie.Domain = "https://api.creatif.app"
			cookie.Path = "/api/v1"
			if os.Getenv("APP_ENV") != "prod" {
				cookie.HttpOnly = false
				cookie.Secure = false
				cookie.Domain = "http://localhost:3000"
				cookie.Path = "/"
			}

			cookie.Value = string(encryptedUser)
			cookie.Expires = time.Now().Add(1 * time.Hour)
			c.SetCookie(cookie)
		})
	}
}
