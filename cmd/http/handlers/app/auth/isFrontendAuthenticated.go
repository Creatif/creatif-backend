package auth

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/auth/isFrontendAuthenticated"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateIsFrontendAuthenticated() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := logger.NewLogBuilder()
		authentication := auth.NewFrontendAuthentication(request.GetAuthenticationCookie(c), l)
		handler := isFrontendAuthenticated.New(authentication, l)

		return request.SendResponse(handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		})
	}
}
