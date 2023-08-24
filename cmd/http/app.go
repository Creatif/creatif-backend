package main

import (
	"creatif/cmd/http/handlers/assignments"
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
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderAccessControlAllowCredentials},
		AllowMethods:     []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	appRoutes(srv.Group("/api/v1/app"))
	declarationRoutes(srv.Group("/api/v1/declarations"))
	assignmentRoutes(srv.Group("/api/v1/assignments"))

	server.StartServer(srv)
}

func appRoutes(group *echo.Group) {
}

func declarationRoutes(group *echo.Group) {
	group.PUT("/node", declarations.CreateNodeHandler())
	group.GET("/node/:id", declarations.GetNodeHandler())
	// TODO test on frontend if query parameters work
	group.GET("/node/map", declarations.GetMapHandler())
}

func assignmentRoutes(group *echo.Group) {
	group.PUT("/node", assignments.AssignNodeHandler())
}
