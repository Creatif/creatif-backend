package project

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/hasProjects"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HasProjectsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := logger.NewLogBuilder()
		handler := hasProjects.New(auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c), l), l)

		return request.SendResponse[interface{}](handler, c, http.StatusOK, l, nil, false)
	}
}
