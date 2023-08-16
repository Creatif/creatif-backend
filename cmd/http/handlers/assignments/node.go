package assignments

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/assignments"
	"creatif/pkg/app/assignments/create"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func AssignNodeHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		b, _ := io.ReadAll(c.Request().Body)

		if model, err := sdk.UnmarshalToStruct[assignments.AssignNodeTextValue](b); err == nil {
			handler := create.New(create.NewCreateNodeModel(model.Name, "text", create.AssignNodeTextModel{
				Name:  model.Name,
				Value: []byte(model.Value),
			}))

			return request.SendResponse[*create.CreateNodeModel](handler, c, http.StatusCreated)
		} else if model, err := sdk.UnmarshalToStruct[assignments.AssignNodeBooleanValue](b); err == nil {
			handler := create.New(create.NewCreateNodeModel(model.Name, "boolean", create.AssignNodeBooleanModel{
				Name:  model.Name,
				Value: model.Value,
			}))

			return request.SendResponse[*create.CreateNodeModel](handler, c, http.StatusCreated)
		}

		return c.JSON(http.StatusInternalServerError, request.ErrorResponse[string]{
			Data: "Internal server error",
		})
	}
}
