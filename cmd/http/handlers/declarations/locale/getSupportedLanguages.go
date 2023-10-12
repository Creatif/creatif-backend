package locale

import (
	"creatif/pkg/app/services/locales"
	"github.com/labstack/echo/v4"
	"net/http"
)

type localeView struct {
	Name  string `json:"name"`
	Alpha string `json:"alpha"`
}

func processStoredLocales(l map[string]map[string]string) []localeView {
	loadedLocales := make([]localeView, len(l))
	i := 0
	for key, lang := range l {
		loadedLocales[i] = localeView{
			Name:  lang["name"],
			Alpha: key,
		}
		i++
	}

	return loadedLocales
}

func GetSupportedLanguageHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		l := locales.Locales()
		if len(l) > 0 {
			return c.JSON(http.StatusOK, processStoredLocales(l))
		}

		return c.JSON(http.StatusOK, processStoredLocales(l))
	}
}
