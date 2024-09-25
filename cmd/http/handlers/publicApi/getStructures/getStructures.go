package getStructures

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	getStructuresService "creatif/pkg/app/services/publicApi/getStructures"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetStructuresHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.GetStructures
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizeGetStructures(model)

		handler := getStructuresService.New(getStructuresService.NewModel(model.ProjectID), auth.NewAnonymousAuthentication())

		return request.SendPublicResponse[getStructuresService.Model](handler, c, http.StatusOK, nil, false)
	}
}
