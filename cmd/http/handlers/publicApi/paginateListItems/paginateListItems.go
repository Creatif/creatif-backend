package paginateListItems

import (
	publicApi2 "creatif/cmd/http/handlers/publicApi"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/publicApi"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/publicApi/paginateListItems"
	"creatif/pkg/app/services/shared/queryProcessor"
	"creatif/pkg/lib/sdk"
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
		model, err := publicApi.SanitizePaginateListItems(model)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"query": "'query' data is invalid and cannot be unpacked",
			})
		}

		handler := paginateListItems.New(paginateListItems.NewModel(
			model.VersionName,
			model.ProjectID,
			model.ListName,
			model.Page,
			model.Limit,
			model.OrderDirection,
			model.OrderBy,
			model.Search,
			model.SanitizedLocales,
			model.SanitizedGroups,
			paginateListItems.Options{
				ValueOnly: model.ResolvedOptions.ValueOnly,
			},
			sdk.Map(model.SanitizedQuery, func(idx int, value publicApi.Query) queryProcessor.Query {
				return queryProcessor.Query{
					Column:   value.Column,
					Value:    value.Value,
					Operator: value.Operator,
				}
			}),
		), auth.NewAnonymousAuthentication())

		return request.SendPublicResponse[paginateListItems.Model](handler, c, http.StatusOK, nil, false)
	}
}
