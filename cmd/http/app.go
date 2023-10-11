package main

import (
	appHandlers "creatif/cmd/http/handlers/app"
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
	releaseAllLocks()
	loadLanguages()

	srv := setupServer()
	srv.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderAccessControlAllowCredentials},
		AllowMethods:     []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	declarationRoutes(srv.Group("/api/v1/declarations"))
	appRoutes(srv.Group("/api/v1/app"))

	server.StartServer(srv)
}

func appRoutes(group *echo.Group) {
	group.PUT("/project", appHandlers.CreateProjectHandler())
}

func declarationRoutes(group *echo.Group) {
	group.GET("/supported-languages", appHandlers.GetSupportedLanguageHandler())
	group.PUT("/variable", declarations.CreateVariableHandler())
	group.POST("/variable/:projectID", declarations.UpdateVariableHandler())
	group.DELETE("/variable/:projectID/:name", declarations.DeleteVariableHandler())
	group.GET("/variable/:projectID/:name", declarations.GetVariableHandler())
	group.GET("/variables/:projectID", declarations.PaginateVariablesHandler())
	group.GET("/variable/value/:projectID/:name", declarations.GetValueHandler())

	group.POST("/map/add/:projectID", declarations.AddToMapHandler())
	group.POST("/map/update/:projectID", declarations.UpdateMapVariableHandler())
	group.DELETE("/map/entry/:projectID/:name/:entryName", declarations.DeleteMapEntry())
	group.PUT("/map/:projectID", declarations.CreateMapHandler())
	group.GET("/map/:projectID/:name", declarations.GetMapHandler())

	group.POST("/structures/:projectID", declarations.GetBatchedStructuresHandler())
}
