package auth

import (
	"context"
	"creatif/cmd"
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/cache"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

func CreateApiAuthSessionHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get(cmd.CreatifApiHeader)
		if apiKey == "" {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		projectId := c.Request().Header.Get(cmd.CreatifProjectIDHeader)
		if projectId == "" {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		var user app.Project
		if res := storage.Gorm().Where("id = ? AND api_key = ?", projectId, apiKey).Select("id").First(&user); res.Error != nil {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		l := logger.NewLogBuilder()
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		id := ksuid.New().String()
		apiKeyProjectId := fmt.Sprintf("%s-%s", apiKey, projectId)

		_, err := cache.Cache().Set(ctx, id, apiKeyProjectId, 5*time.Minute).Result()
		if err != nil {
			l.Add("cache.fail", err.Error())
			l.Flush("error")
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		return c.JSON(http.StatusOK, id)
	}
}
