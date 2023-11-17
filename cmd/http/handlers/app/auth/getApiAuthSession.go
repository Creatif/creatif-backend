package auth

import (
	"context"
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
		var model app.GetApiAuthSession
		if err := c.Bind(&model); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		s := strings.Split(model.Session, "-")
		if len(s) != 2 {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		l := logger.NewLogBuilder()
		ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		key, err := cache.Cache().Get(ctx, s[0]).Result()
		if err != nil {
			l.Add("cache.fail", err.Error())
			l.Flush("error")
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		if key != s[1] {
			return c.JSON(http.StatusForbidden, "Unauthenticated")
		}

		return c.NoContent(http.StatusOK)
	}
}
