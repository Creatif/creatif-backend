package stats

import (
	"creatif/cmd/http/request"
	"creatif/cmd/http/request/app"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/stats/dashboard"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetDashboardStatsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		var model app.GetDashboardStats
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		model = app.SanitizeGetDashboardStats(model)

		handler := dashboard.New(dashboard.NewModel(model.ProjectID), auth.NewApiAuthentication(request.GetAuthenticationCookie(c)))

		return request.SendResponse[dashboard.Model](handler, c, http.StatusOK, nil, false)
	}
}
