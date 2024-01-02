package maps

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/maps"
	"creatif/pkg/app/auth"
	mapCreate2 "creatif/pkg/app/services/maps/mapCreate"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateMapHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model maps.CreateMap
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = maps.SanitizeMapModel(model)

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		serviceEntries := make([]mapCreate2.VariableModel, 0)
		for _, entry := range model.Variables {
			serviceEntries = append(serviceEntries, mapCreate2.VariableModel{
				Name:      entry.Name,
				Metadata:  []byte(entry.Metadata),
				Locale:    entry.Locale,
				Groups:    entry.Groups,
				Behaviour: entry.Behaviour,
				Value:     []byte(entry.Value),
			})
		}

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := mapCreate2.New(mapCreate2.NewModel(model.ProjectID, model.Name, serviceEntries), authentication, l)

		return request.SendResponse[mapCreate2.Model](handler, c, http.StatusCreated, l, func(c echo.Context, model interface{}) error {
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
