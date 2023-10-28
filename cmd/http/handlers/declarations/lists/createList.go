package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/createList"
	"creatif/pkg/lib/logger"
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
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		l := logger.NewLogBuilder()
		handler := createList.New(createList.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			sdk.Map(model.Variables, func(idx int, value lists.CreateListVariable) createList.Variable {
				return createList.Variable{
					Name:      value.Name,
					Metadata:  []byte(value.Metadata),
					Groups:    value.Groups,
					Behaviour: value.Behaviour,
					Value:     []byte(value.Value),
				}
			}),
		))

		return request.SendResponse[createList.Model](handler, c, http.StatusCreated, l)
	}
}
