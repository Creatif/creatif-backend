package paginateMapItems

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
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

		versionName := c.Request().Header.Get(publicApi2.CreatifVersionHeader)
		model.VersionName = versionName
		model = publicApi.SanitizePaginateMapItems(model)

		l := logger.NewLogBuilder()
		handler := paginateMapItems.New(paginateMapItems.NewModel(
			model.VersionName,
			model.ProjectID,
			model.ListName,
			model.Page,
			model.OrderDirection,
			model.OrderBy,
			model.Search,
			model.SanitizedLocales,
			model.SanitizedGroups,
			paginateMapItems.Options{
				ValueOnly: model.ResolvedOptions.ValueOnly,
			},
		), auth.NewAnonymousAuthentication(), l)

		return request.SendPublicResponse[paginateMapItems.Model](handler, c, http.StatusOK, l, nil, false)
	}
}
