package project

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/createProject"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateProjectHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.CreateProject
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeProject(model)
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := createProject.New(createProject.NewModel(model.Name), authentication)

		return request.SendResponse[createProject.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
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
