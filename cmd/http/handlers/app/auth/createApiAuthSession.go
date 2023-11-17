package auth

import (
	"context"
	"creatif/pkg/lib/cache"
	"creatif/pkg/lib/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
	"net/http"
	"time"
)

func CreateApiAuthSessionHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := logger.NewLogBuilder()
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		id := ksuid.New().String()
		value := ksuid.New().String()

		_, err := cache.Cache().Set(ctx, id, value, 5*time.Minute).Result()
		if err != nil {
			l.Add("cache.fail", err.Error())
			l.Flush("error")
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		return c.JSON(http.StatusOK, fmt.Sprintf("%s-%s", id, value))
	}
}
