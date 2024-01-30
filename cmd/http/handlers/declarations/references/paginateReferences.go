package references

import (
	"creatif/cmd"
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/declarations/references"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/shared/paginateReferences"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/sdk"
	"github.com/labstack/echo/v4"
	"net/http"
)

func PaginateReferencesHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model references.PaginateReferences
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = references.SanitizePaginateReferences(model)

		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)

		l := logger.NewLogBuilder()
		authentication := auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), projectId, apiKey, l)
		handler := paginateReferences.New(paginateReferences.NewModel(
			model.ProjectID,
			model.ParentID,
			model.ChildID,
			model.ParentStructureID,
			model.ChildStructureID,
			model.RelationshipType,
			model.StructureType,
			model.SanitizedLocales,
			model.OrderBy,
			model.Search,
			model.OrderDirection,
			model.Limit,
			model.Page,
			model.SanitizedGroups,
			sdk.ParseFilters(model.Filters),
			model.Behaviour,
			model.SanitizedFields,
		), authentication, l)

		return request.SendResponse[paginateReferences.Model](handler, c, http.StatusOK, l, func(c echo.Context, model interface{}) error {
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
