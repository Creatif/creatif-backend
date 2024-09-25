package paginateListItems

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/paginateListItems"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateListItemsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model publicApi.PaginateListItems
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		versionName := c.Request().Header.Get(publicApi2.CreatifVersionHeader)
		model.VersionName = versionName
		model = publicApi.SanitizePaginateListItems(model)

		handler := paginateListItems.New(paginateListItems.NewModel(
			model.VersionName,
			model.ProjectID,
			model.ListName,
			model.Page,
			model.OrderDirection,
			model.OrderBy,
			model.Search,
			model.SanitizedLocales,
			model.SanitizedGroups,
			paginateListItems.Options{
				ValueOnly: model.ResolvedOptions.ValueOnly,
			},
		), auth.NewAnonymousAuthentication())

		return request.SendPublicResponse[paginateListItems.Model](handler, c, http.StatusOK, nil, false)
	}
}
