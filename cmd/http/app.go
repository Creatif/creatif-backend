package main

import (
	"creatif/cmd/http/handlers/declarations"
	"creatif/cmd/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func app() {
	loadEnv()
	runLogger()
	runAssets()
	runDb()

	srv := setupServer()
	srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderAccessControlAllowCredentials},
		AllowMethods:     []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	declarationRoutes(srv.Group("/api/v1/declarations"))

	server.StartServer(srv)
}

func declarationRoutes(group *echo.Group) {
	group.PUT("/variable", declarations.CreateVariableHandler())
	group.POST("/variable", declarations.UpdateVariableHandler())
	group.DELETE("/variable/:name", declarations.DeleteVariableHandler())
	group.PUT("/map", declarations.CreateMapHandler())
	group.GET("/variable/:name", declarations.GetVariableHandler())
	group.GET("/variables", declarations.PaginateVariablesHandler())
	group.GET("/map/:name", declarations.GetMapHandler())
	group.GET("/variable/value/:name", declarations.GetValueHandler())
	group.POST("/structures", declarations.GetBatchedStructuresHandler())
}
