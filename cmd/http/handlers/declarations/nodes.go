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

		handler := create.New(create.NewCreateNodeModel(model.Name, model.Type, model.Behaviour, model.Groups, []byte(model.Metadata), func() map[string]create.NodeValidation {
			requestValidation := model.Validation
			modelValidation := make(map[string]create.NodeValidation)

			for key, val := range requestValidation {
				validation := create.NodeValidation{
					Required:    val.Required,
					Length:      val.Length,
					ExactValue:  val.ExactValue,
					ExactValues: val.ExactValues,
					IsDate:      val.IsDate,
				}

				modelValidation[key] = validation
			}

			return modelValidation
		}()))

		return request.SendResponse[create.CreateNodeModel](handler, c, http.StatusCreated)
	}
}
