package removeVersion

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publishing/removeVersion"
	"creatif/pkg/app/auth"
	removeVersionService "creatif/pkg/app/services/publishing/removeVersion"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RemoveVersionHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model removeVersion.RemoveVersion
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = removeVersion.SanitizeRemoveVersion(model)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l)
		handler := removeVersionService.New(removeVersionService.NewModel(model.ProjectID, model.ID), authentication, l)

		return request.SendResponse[removeVersionService.Model](handler, c, http.StatusCreated, l, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
