package auth

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LoginApiCheckHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		cookie := request.GetApiAuthenticationCookie(c)
		if cookie == "" {
			fmt.Println("no cookie in loginApiCheck")
			return c.NoContent(http.StatusForbidden)
		}

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		fmt.Println(apiKey, projectId)

		l := logger.NewLogBuilder()
		a := auth.NewApiAuthentication(cookie, projectId, apiKey, l)
		if err := a.Authenticate(); err != nil {
			return c.NoContent(http.StatusForbidden)
		}

		if a.ShouldRefresh() {
			session, err := a.Refresh()
			if err != nil {
				return c.NoContent(http.StatusForbidden)
			}

			c.SetCookie(request.EncryptAuthenticationCookie(session))
		}

		l.Flush("")

		return c.NoContent(http.StatusOK)
	}
}
