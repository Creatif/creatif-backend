package project

import (
	"creatif/cmd/http/request"
	"creatif/pkg/app/auth"
	"creatif/pkg/app/services/projects/hasProjects"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HasProjectsHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		handler := hasProjects.New(auth.NewApiAuthentication(request.GetApiAuthenticationCookie(c)))

		return request.SendResponse[interface{}](handler, c, http.StatusOK, nil, false)
	}
}
