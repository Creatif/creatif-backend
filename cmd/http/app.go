package main

import (
	appHandlers "creatif/cmd/http/handlers/app"
	"creatif/cmd/http/handlers/declarations/combined"
	"creatif/cmd/http/handlers/declarations/lists"
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
	group.POST("/structures/:projectID/:locale", combined.GetBatchedStructuresHandler())

	group.PUT("/list/:projectID", lists.CreateListHandler())
	group.PUT("/list/append/:projectID", lists.AppendToListHandler())
	group.DELETE("/list/:projectID/:name/:locale", lists.DeleteListHandler())
	group.DELETE("/list/item-id/:projectID/:name/:itemID/:locale", lists.DeleteListItemByIDHandler())
	group.POST("/list/range/:projectID/:name/:locale", lists.DeleteRangeByIDHandler())
	group.GET("/lists/:projectID/:name/:locale", lists.PaginateListItemsHandler())
	group.GET("/lists/query-id/:projectID/:name/:locale/:id", lists.QueryListByIDHandler())
	group.POST("/lists/:projectID/:name/:itemName/:locale", lists.ReplaceListItemHandler())
	group.POST("/lists/switch-id/:projectID/:name/:locale/:source/:destination", lists.SwitchByIDHandler())
	group.POST("/list/update/:projectID/:name/:locale", lists.UpdateListHandler())
	group.POST("/list/update-item-by-id/:projectID/:name/:itemID/:locale", lists.UpdateListItemByIDHandler())

	group.POST("/map/add/:projectID/:locale", maps.AddToMapHandler())
	group.POST("/map/update/:projectID/:locale", maps.UpdateMapVariableHandler())
	group.DELETE("/map/entry/:projectID/:name/:entryName/:locale", maps.DeleteMapEntry())
	group.DELETE("/map/:projectID/:name/:locale", maps.DeleteMap())
	group.PUT("/map/:projectID/:locale", maps.CreateMapHandler())
	group.GET("/map/:projectID/:name/:locale", maps.GetMapHandler())

	group.GET("/supported-locales", locale.GetSupportedLocalesHandler())

	group.PUT("/variable", variables.CreateVariableHandler())
	group.POST("/variable/:projectID", variables.UpdateVariableHandler())
	group.DELETE("/variable/:projectID", variables.DeleteVariableHandler())
	group.GET("/variable/:projectID/:name/:locale", variables.GetVariableHandler())
	group.GET("/variables/:projectID/:locale", variables.PaginateVariablesHandler())
	group.GET("/variable/value/:projectID/:name/:locale", variables.GetValueHandler())
}
