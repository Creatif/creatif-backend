package assignments

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/assignments"
	"creatif/pkg/app/assignments/mapCreate"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func AssignMapValueHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		b, _ := io.ReadAll(c.Request().Body)

		if model, err := sdk.UnmarshalToStruct[assignments.AssignNodeTextValue](b); err == nil {
			model = assignments.SanitizeTextValue(model)
			handler := mapCreate.New(mapCreate.NewAssignValueModel(model.Name, []byte(model.Value)))

			return request.SendResponse[*mapCreate.AssignValueModel](handler, c, http.StatusCreated)
		}

		return c.JSON(http.StatusInternalServerError, request.ErrorResponse[string]{
			Data: "Internal server error",
		})
	}
}
