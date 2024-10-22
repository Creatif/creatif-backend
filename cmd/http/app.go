package main

import (
	authHandlers "creatif/cmd/http/handlers/app/auth"
	"creatif/cmd/http/handlers/app/groups"
	appHandlers "creatif/cmd/http/handlers/app/project"
	"creatif/cmd/http/handlers/app/structures"
	"creatif/cmd/http/handlers/declarations/files"
	"creatif/cmd/http/handlers/declarations/lists"
	"creatif/cmd/http/handlers/declarations/locale"
	"creatif/cmd/http/handlers/declarations/maps"
	"creatif/cmd/http/handlers/declarations/references"
	"creatif/cmd/http/handlers/publicApi/getFile"
	"creatif/cmd/http/handlers/publicApi/getListItemByID"
	"creatif/cmd/http/handlers/publicApi/getListItemsByName"
	"creatif/cmd/http/handlers/publicApi/getMany"
	"creatif/cmd/http/handlers/publicApi/getMapItemByID"
	"creatif/cmd/http/handlers/publicApi/getMapItemByName"
	"creatif/cmd/http/handlers/publicApi/getStructures"
	"creatif/cmd/http/handlers/publicApi/getVersions"
	"creatif/cmd/http/handlers/publicApi/paginateListItems"
	"creatif/cmd/http/handlers/publicApi/paginateMapItems"
	"creatif/cmd/http/handlers/publishing/publish"
	"creatif/cmd/http/handlers/publishing/removeVersion"
	"creatif/cmd/http/handlers/publishing/toggleProduction"
	"creatif/cmd/server"
	"creatif/pkg/app/services/events"
	"creatif/pkg/app/services/publicApi/publicApiError"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
)

func app() {
	loadEnv()
	runAssets()
	runPublic()
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

	createDatabase()

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
			"Creatif-Version",
		},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT, echo.DELETE},
	}))

	healthRoutes(srv.Group("/api/v1/health"))
	declarationRoutes(srv.Group("/api/v1/declarations"))
	appRoutes(srv.Group("/api/v1/app"))
	publishingRoutes(srv.Group("/api/v1/publishing"))
	publicRoutes(srv.Group("/api/v1/public"))
	appFiles(srv.Group("/api/v1/files"))

	events.RunEvents()
	server.StartServer(srv)
}

func healthRoutes(group *echo.Group) {
	group.GET("/full-health", func(c echo.Context) error {
		db, err := storage.Gorm().DB()
		if err != nil {
			fmt.Println("Cannot get db", err)
			return c.String(http.StatusServiceUnavailable, err.Error())
		}

		if err := db.Ping(); err != nil {
			fmt.Println("Cannot ping", err)
			return c.String(http.StatusServiceUnavailable, err.Error())
		}

		return c.String(http.StatusOK, "HEALTHY")
	})
}

func appRoutes(group *echo.Group) {
	group.GET("/supported-locales", locale.GetSupportedLocalesHandler())

	group.PUT("/project", appHandlers.CreateProjectHandler())
	group.POST("/project/structure/truncate/:projectId", structures.TruncateStructureHandler())
	group.POST("/project/structure/remove/:projectId", structures.RemoveStructureHandler())
	group.GET("/project/exists", appHandlers.HasProjectsHandler())
	group.POST("/project/metadata/:projectId", structures.GetStructureMetadataHandler())
	group.GET("/projects", appHandlers.PaginateProjectsHandler())
	group.GET("/project/single/:id", appHandlers.GetProjectHandler())
	group.PUT("/groups/:projectId", groups.AddGroupsHandler())
	group.GET("/groups/:projectId", groups.GetGroupsHandler())

	group.PUT("/auth/admin/create", authHandlers.CreateAdminHandler())
	group.GET("/auth/admin/exists", authHandlers.AdminExistsHandler())
	group.POST("/auth/login", authHandlers.LoginHandler())

	group.POST("/auth/logout", authHandlers.LogoutApiHandler())
}

func appFiles(group *echo.Group) {
	group.GET("/file/:projectID/:id", files.GetFileHandler())
}

func declarationRoutes(group *echo.Group) {
	group.PUT("/list/:projectID", lists.CreateListHandler())
	group.PUT("/list/add/:projectID", lists.AddToListHandler())
	group.DELETE("/list/:projectID/:name", lists.DeleteListHandler())
	group.POST("/list/item-id/:projectID", lists.DeleteListItemByIDHandler())
	group.GET("/list/groups/:projectID/:name/:itemId", lists.GetListGroupsHandler())
	group.POST("/list/range/:projectID/:name", lists.DeleteRangeByIDHandler())
	group.GET("/lists/items/:projectID/:name", lists.PaginateListItemsHandler())
	group.GET("/lists/:projectID", lists.PaginateListsHandler())
	group.GET("/list/query-id/:projectID/:name/:itemId", lists.QueryListByIDHandler())
	group.POST("/list/rearrange/:projectID/:name/:source/:destination", lists.SwitchByIDHandler())
	group.POST("/list/update/:projectID/:name", lists.UpdateListHandler())
	group.POST("/list/update-item-by-id/:projectID/:name/:itemID", lists.UpdateListItemByIDHandler())

	group.PUT("/map/add/:projectID/:locale", maps.AddToMapHandler())
	group.POST("/map/update/:projectID/:name/:itemId", maps.UpdateMapVariableHandler())
	group.DELETE("/map/entry/:projectID/:name/:variableName", maps.DeleteMapEntry())
	group.DELETE("/map/:projectID/:name/:locale", maps.DeleteMap())
	group.GET("/map/query-id/:projectID/:name/:itemId", maps.QueryMapVariableHandler())
	group.PUT("/map/:projectID", maps.CreateMapHandler())
	group.POST("/map/rearrange/:projectID/:name/:source/:destination", maps.SwitchByIDHandler())
	group.GET("/map/groups/:projectID/:name/:itemId", maps.GetMapGroupsHandler())
	group.POST("/map/range/:projectID/:name", maps.DeleteRange())
	group.GET("/maps/items/:projectID/:name", maps.PaginateMapVariables())
	group.GET("/maps/:projectID", maps.PaginateMapsHandler())
	group.GET("/map/:projectID/:name", maps.GetMapHandler())

	group.GET("/references/:projectID/:parentID/:childID/:structureType/:relationshipType/:childStructureId/:parentStructureId", references.PaginateReferencesHandler())
}

func publishingRoutes(group *echo.Group) {
	group.PUT("/publish/:projectId", publish.PublishHandler())
	group.DELETE("/publish/version/:projectId/:id", removeVersion.RemoveVersionHandler())
	group.POST("/publish/toggle-production/:projectId/:id", toggleProduction.ToggleProductionHandler())
}

func publicRoutes(group *echo.Group) {
	group.GET("/:projectId/versions", getVersions.GetVersionsHandler())
	group.GET("/:projectId/many", getMany.GetManyHandler())
	group.GET("/:projectId/list/:structureName/:name", getListItemsByName.GetListItemsByNameHandler())
	group.GET("/:projectId/map/:structureName/:name", getMapItemByName.GetMapItemByNameHandler())
	group.GET("/:projectId/structures", getStructures.GetStructuresHandler())
	group.GET("/:projectId/list/id/:id", getListItemByID.GetListItemByIDHandler())
	group.GET("/:projectId/map/id/:id", getMapItemByID.GetMapItemByIDHandler())
	group.GET("/:projectId/lists/:name", paginateListItems.PaginateListItemsHandler())
	group.GET("/:projectId/maps/:name", paginateMapItems.PaginateMapItemsHandler())
	group.GET("/:projectId/file/:version/:id", getFile.GetFileHandler())

	group.Any("/:projectId/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"call": "unknown",
			"messages": map[string]string{
				"notFound": "This route does not exist",
			},
			"status": publicApiError.NotFoundError,
		})
	})
}
