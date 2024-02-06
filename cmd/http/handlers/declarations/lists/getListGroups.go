package lists

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/getListGroups"
	"creatif/pkg/lib/logger"
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

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		a := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := getListGroups.New(getListGroups.NewModel(
			model.Name,
			model.ItemID,
			model.ProjectID,
		), a, l)

		return request.SendResponse(handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
