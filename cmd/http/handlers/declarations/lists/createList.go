package lists

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/lists/createList"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateListHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.CreateList
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeCreateList(model)

		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c))
		handler := createList.New(createList.NewModel(
			model.ProjectID,
			model.Name,
			sdk.Map(model.Variables, func(idx int, value lists.CreateListVariable) createList.Variable {
				return createList.Variable{
					Name:      value.Name,
					Metadata:  []byte(value.Metadata),
					Groups:    value.Groups,
					Locale:    value.Locale,
					Behaviour: value.Behaviour,
					Value:     []byte(value.Value),
				}
			}),
		), authentication)

		return request.SendResponse[createList.Model](handler, c, http.StatusCreated, func(c echo.Context, model interface{}) error {
			if authentication.ShouldRefresh() {
				session, err := authentication.Refresh()
				if err != nil {
					return err
				}

				c.SetCookie(request.EncryptAuthenticationCookie(session))
			}

			return nil
		}, model.GracefulFail)
	}
}
