package handlers

import (
	"creatif/cmd/http/request"
	projectCreate "creatif/pkg/app/projects/create"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateProjectHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model request.CreateProject
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		model = request.SanitizeProject(model)

		handler := projectCreate.New(projectCreate.NewCreateProjectModel(model.Name))

		return request.SendResponse[projectCreate.CreateProjectModel](handler, c, http.StatusCreated)
	}
}
