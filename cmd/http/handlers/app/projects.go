package app

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	projectCreate "creatif/pkg/app/projects/create"
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

		handler := projectCreate.New(projectCreate.NewCreateProjectModel(model.Name))

		return request.SendResponse[projectCreate.CreateProjectModel](handler, c, http.StatusCreated)
	}
}
