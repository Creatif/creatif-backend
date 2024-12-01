package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/queryListByID"
	"github.com/labstack/echo/v4"
	"net/http"
)

func QueryListByIDHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.QueryListByID
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeQueryListByID(model)
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := queryListByID.New(queryListByID.NewModel(
			model.ProjectID,
			model.Name,
			model.ItemID,
			model.ConnectionReplaceMethod,
		), authentication)

		return request.SendResponse[queryListByID.Model](handler, c, http.StatusOK, func(c echo.Context, model interface{}) error {
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
