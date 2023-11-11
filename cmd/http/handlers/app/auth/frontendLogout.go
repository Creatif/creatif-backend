package auth

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"time"
)

func CreateFrontendLogout() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := logger.NewLogBuilder()
		authentication := auth.NewFrontendAuthentication(request.GetAuthenticationCookie(c), l)

		authentication.Logout(func() {
			cookie := new(http.Cookie)
			cookie.Name = "authentication"
			cookie.HttpOnly = true
			cookie.Secure = true
			cookie.Domain = "https://api.creatif.app"
			cookie.Path = "/api/v1"
			if os.Getenv("APP_ENV") != "prod" {
				cookie.HttpOnly = true
				cookie.Secure = true
				cookie.Domain = "http://localhost"
				cookie.Path = "/"
			}

			cookie.Value = ""
			cookie.Expires = time.Unix(0, 0)

			c.SetCookie(cookie)
		})

		return c.JSON(http.StatusOK, nil)
	}
}
