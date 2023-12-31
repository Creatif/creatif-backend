package auth

import (
	"creatif/cmd/http/request"
	"github.com/labstack/echo/v4"
	"net/http"
)

func LogoutApiHandler() func(e echo.Context) error {
	return func(c echo.Context) error {

		c.SetCookie(request.RemoveApiAuthenticationCookie())

		return c.JSON(http.StatusOK, nil)
	}
}
