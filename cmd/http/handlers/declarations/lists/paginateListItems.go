package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/paginateListItems"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateListItemsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.PaginateListItems
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizePaginateListItems(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := paginateListItems.New(paginateListItems.NewModel(
			model.ProjectID,
			model.SanitizedLocales,
			model.ListName,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
			model.Behaviour,
			model.SanitizedFields,
		), authentication)

		return request.SendResponse[paginateListItems.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
