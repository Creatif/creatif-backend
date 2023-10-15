package lists

import (
	declarations2 "creatif/cmd/http/handlers/declarations"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/lists"
	"creatif/pkg/app/services/lists/replaceListItem"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ReplaceListItemHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model lists.ReplaceListItem
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = lists.SanitizeReplaceListItem(model)
		if model.Locale == "" {
			model.Locale = declarations2.DefaultLocale
		}

		handler := replaceListItem.New(replaceListItem.NewModel(
			model.ProjectID,
			model.Locale,
			model.Name,
			model.ItemName,
			replaceListItem.Variable{
				Name:      model.Variable.Name,
				Metadata:  []byte(model.Variable.Metadata),
				Groups:    model.Variable.Groups,
				Behaviour: model.Variable.Behaviour,
				Value:     []byte(model.Variable.Value),
			},
		))

		return request.SendResponse[replaceListItem.Model](handler, c, http.StatusCreated)
	}
}