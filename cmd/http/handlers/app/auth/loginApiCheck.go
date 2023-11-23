package auth

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
)

func LoginApiCheckHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		cookie := request.GetApiAuthenticationCookie(c)

		if cookie == "" {
			return c.NoContent(403)
		}

		l := logger.NewLogBuilder()
		a := auth.NewApiAuthentication(cookie, l)
		if err := a.Authenticate(); err != nil {
			return c.NoContent(403)
		}

		if a.ShouldRefresh() {
			newToken, err := a.Refresh()
			if err != nil {
				return c.NoContent(403)
			}

			newAuthCookie := request.EncryptApiAuthenticationCookie(newToken)
			c.SetCookie(newAuthCookie)

			return c.NoContent(200)
		}

		return c.NoContent(200)
	}
}
