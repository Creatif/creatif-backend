package auth

import (
	"context"
	"creatif/cmd"
	"creatif/cmd/http/request/app"
	"creatif/pkg/lib/cache"
	"creatif/pkg/lib/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

func GetApiAuthSession() func(e echo.Context) error {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		if apiKey == "" {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)
		if apiKey == "" {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		var model app.GetApiAuthSession
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		l := logger.NewLogBuilder()
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		key, err := cache.Cache().Get(ctx, model.Session).Result()
		if err != nil {
			l.Add("cache.fail", err.Error())
			l.Flush("error")
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		split := strings.Split(key, "-")
		if len(split) != 2 {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		cacheApiKey := split[0]
		cacheProjectId := split[1]

		if cacheApiKey != apiKey && cacheProjectId != projectId {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		return c.NoContent(http.StatusOK)
	}
}
