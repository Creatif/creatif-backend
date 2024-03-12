package paginateMapItems

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/paginateMapItems"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateMapItemsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.PaginateMapItems
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = publicApi.SanitizePaginateMapItems(model)

		l := logger.NewLogBuilder()
		handler := paginateMapItems.New(paginateMapItems.NewModel(
			model.ProjectID,
			model.ListName,
			model.Page,
			model.OrderDirection,
			model.OrderBy,
			model.Search,
			model.SanitizedLocales,
			model.SanitizedGroups,
		), auth.NewAnonymousAuthentication(), l)

		return request.SendResponse[paginateMapItems.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
