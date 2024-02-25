package publish

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	publishRequest "creatif/cmd/http/request/publishing/publish"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publishing/publish"
	"creatif/pkg/lib/logger"
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

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := publish.New(publish.NewModel(model.ProjectID, model.Name), authentication, l)

		return request.SendResponse[publish.Model](handler, c, http.StatusCreated, l, func(c echo.Context, model interface{}) error {
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
