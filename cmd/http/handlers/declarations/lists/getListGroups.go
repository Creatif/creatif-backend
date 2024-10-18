package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/getListGroups"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetListGroupsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.GetListGroups
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeGetListGroups(model)

		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := getListGroups.New(getListGroups.NewModel(
			model.Name,
			model.ItemID,
			model.ProjectID,
		), a)

		return request.SendResponse[getListGroups.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
			if a.ShouldRefresh() {
				session, err := a.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, false)
	}
}
