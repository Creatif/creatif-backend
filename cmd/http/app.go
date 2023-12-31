package main

import (
	"creatif/cmd"
	authHandlers "creatif/cmd/http/handlers/app/auth"
	appHandlers "creatif/cmd/http/handlers/app/project"
	"creatif/cmd/http/handlers/declarations/combined"
	"creatif/cmd/http/handlers/declarations/lists"
	"creatif/cmd/http/handlers/declarations/locale"
	"creatif/cmd/http/handlers/declarations/maps"
	"creatif/cmd/http/handlers/declarations/variables"
	"creatif/cmd/server"
	"creatif/pkg/lib/cache"
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
	if err := cache.NewCache(); err != nil {
		log.Fatalln(err)
	}

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
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderCookie,
			echo.HeaderAccessControlAllowCredentials,
			cmd.CreatifApiHeader,
			cmd.CreatifProjectIDHeader,
		},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	declarationRoutes(srv.Group("/api/v1/declarations"))
	appRoutes(srv.Group("/api/v1/app"))

	server.StartServer(srv)
}

func appRoutes(group *echo.Group) {
	group.GET("/supported-locales", locale.GetSupportedLocalesHandler())

	group.PUT("/project", appHandlers.CreateProjectHandler())
	group.GET("/project-metadata", appHandlers.GetProjectMetadataHandler())
	group.GET("/projects", appHandlers.PaginateProjectsHandler())
	group.GET("/project/:id", appHandlers.GetProjectHandler())

	group.PUT("/auth/register/email", authHandlers.CreateRegisterEmailHandler())
	group.POST("/auth/login/email", authHandlers.CreateLoginEmailHandler())
	group.POST("/auth/login/api", authHandlers.CreateLoginApiHandler())
	group.POST("/auth/frontend-authenticated", authHandlers.CreateIsFrontendAuthenticated())
	group.POST("/auth/frontend-logout", authHandlers.CreateFrontendLogout())

	group.POST("/auth/api-auth-session", authHandlers.CreateApiAuthSessionHandler())
	group.POST("/auth/api-check", authHandlers.LoginApiCheckHandler())
	group.POST("/auth/logout", authHandlers.LogoutApiHandler())
	group.GET("/auth/api-auth-session/:session", authHandlers.GetApiAuthSession())
}

func declarationRoutes(group *echo.Group) {
	group.POST("/structures/:projectID/:locale", combined.GetBatchedStructuresHandler())

	group.PUT("/list/:projectID/:locale", lists.CreateListHandler())
	group.PUT("/list/append/:projectID", lists.AppendToListHandler())
	group.DELETE("/list/:projectID/:name", lists.DeleteListHandler())
	group.POST("/list/item-id/:projectID", lists.DeleteListItemByIDHandler())
	group.GET("/list/groups/:projectID/:name", lists.GetListGroupsHandler())
	group.POST("/list/range/:projectID/:name", lists.DeleteRangeByIDHandler())
	group.GET("/lists/:projectID/:name", lists.PaginateListItemsHandler())
	group.GET("/list/query-id/:projectID/:name/:itemId", lists.QueryListByIDHandler())
	group.POST("/lists/:projectID/:name/:itemName", lists.ReplaceListItemHandler())
	group.POST("/lists/rearrange/:projectID/:name/:source/:destination", lists.SwitchByIDHandler())
	group.POST("/list/update/:projectID/:name", lists.UpdateListHandler())
	group.POST("/list/update-item-by-id/:projectID/:name/:itemID", lists.UpdateListItemByIDHandler())

	group.POST("/map/add/:projectID/:locale", maps.AddToMapHandler())
	group.POST("/map/update/:projectID/:locale", maps.UpdateMapVariableHandler())
	group.DELETE("/map/entry/:projectID/:name/:entryName/:locale", maps.DeleteMapEntry())
	group.DELETE("/map/:projectID/:name/:locale", maps.DeleteMap())
	group.PUT("/map/:projectID/:locale", maps.CreateMapHandler())
	group.GET("/map/:projectID/:name/:locale", maps.GetMapHandler())

	group.PUT("/variable/:projectID/:locale", variables.CreateVariableHandler())
	group.POST("/variable/:projectID", variables.UpdateVariableHandler())
	group.DELETE("/variable/:projectID/:locale/:name", variables.DeleteVariableHandler())
	group.GET("/variable/:projectID/:name/:locale", variables.GetVariableHandler())
	group.GET("/variable/groups/:projectID/:name", variables.GetVariableGroupsHandler())
	group.GET("/variables/:projectID", variables.PaginateVariablesHandler())
	group.GET("/variable/value/:projectID/:name/:locale", variables.GetValueHandler())
}
