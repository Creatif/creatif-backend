package main

import (
	appHandlers "creatif/cmd/http/handlers/app"
	"creatif/cmd/http/handlers/declarations/combined"
	"creatif/cmd/http/handlers/declarations/locale"
	"creatif/cmd/http/handlers/declarations/maps"
	"creatif/cmd/http/handlers/declarations/variables"
	"creatif/cmd/server"
	"creatif/pkg/lib/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

func app() {
	loadEnv()
	runLogger()
	runAssets()
	runDb()
	if err := releaseAllLocks(); err != nil {
		sqlDB, err := storage.SQLDB()
		if err != nil {
			log.Fatalln("Unable to get storage.SQLDB()", err)
		}

		if err := sqlDB.Close(); err != nil {
			log.Fatalln("Unable to disconnect from the database", err)
		}
	}

	if err := loadLocales(); err != nil {
		sqlDB, err := storage.SQLDB()
		if err != nil {
			log.Fatalln("Unable to get storage.SQLDB()", err)
		}

		if err := sqlDB.Close(); err != nil {
			log.Fatalln("Unable to disconnect from the database", err)
		}
	}

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
	group.GET("/supported-locales", locale.GetSupportedLanguageHandler())
	group.PUT("/variable", variables.CreateVariableHandler())
	group.POST("/variable/:projectID", variables.UpdateVariableHandler())
	group.DELETE("/variable/:projectID/:name", variables.DeleteVariableHandler())
	group.GET("/variable/:projectID/:name", variables.GetVariableHandler())
	group.GET("/variables/:projectID", variables.PaginateVariablesHandler())
	group.GET("/variable/value/:projectID/:name", variables.GetValueHandler())

	group.POST("/map/add/:projectID", maps.AddToMapHandler())
	group.POST("/map/update/:projectID", maps.UpdateMapVariableHandler())
	group.DELETE("/map/entry/:projectID/:name/:entryName", maps.DeleteMapEntry())
	group.PUT("/map/:projectID", maps.CreateMapHandler())
	group.GET("/map/:projectID/:name", maps.GetMapHandler())

	group.POST("/structures/:projectID", combined.GetBatchedStructuresHandler())
}
