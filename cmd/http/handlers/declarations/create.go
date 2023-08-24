package declarations

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations"
	create "creatif/pkg/app/declarations/create"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateNodeHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model declarations.CreateNode
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = declarations.SanitizeNode(model)

		handler := create.New(create.NewCreateNodeModel(model.Name, model.Type, model.Behaviour, model.Groups, []byte(model.Metadata), func() create.NodeValidation {
			requestValidation := model.Validation

			return create.NodeValidation{
				Required: requestValidation.Required,
				Length: create.ValidationLength{
					Min:   requestValidation.Length.Min,
					Max:   requestValidation.Length.Max,
					Exact: requestValidation.Length.Exact,
				},
				ExactValue:  requestValidation.ExactValue,
				ExactValues: requestValidation.ExactValues,
				IsDate:      requestValidation.IsDate,
			}
		}()))

		return request.SendResponse[create.CreateNodeModel](handler, c, http.StatusCreated)
	}
}
