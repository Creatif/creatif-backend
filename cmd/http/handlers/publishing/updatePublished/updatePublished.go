package updatePublished

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publishing/updatePublished"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publishing/updateVersion"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PublishUpdateHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model updatePublished.UpdatePublished
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = updatePublished.SanitizeUpdatePublished(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := updateVersion.New(updateVersion.NewModel(model.ProjectID, model.Name), authentication)

		return request.SendResponse[updateVersion.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
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
