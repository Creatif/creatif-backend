package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/paginateLists"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateListsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.PaginateLists
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizePaginateLists(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := paginateLists.New(paginateLists.NewModel(
			model.ProjectID,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
		), authentication)

		return request.SendResponse[paginateLists.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
