package publish

import (
	"creatif/cmd/http/request"
	publishRequest "creatif/cmd/http/request/publishing/publish"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publishing/publish"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PublishHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publishRequest.Publish
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publishRequest.SanitizePublish(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := publish.New(publish.NewModel(model.ProjectID, model.Name), authentication)

		return request.SendResponse[publish.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
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
