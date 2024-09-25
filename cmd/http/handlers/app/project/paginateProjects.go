package project

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/paginateProjects"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateProjectsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.PaginateProjects
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizePaginateProjects(model)

		handler := paginateProjects.New(paginateProjects.NewModel(
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
		), auth.NewFrontendAuthentication(request.GetAuthenticationCookie(c)))

		return request.SendResponse[paginateProjects.Model](handler, c, http.StatusOK, nil, false)
	}
}
