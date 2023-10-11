package declarations

import (
	"creatif/pkg/app/services/languages"
	"github.com/labstack/echo/v4"
	"net/http"
)

type languageView struct {
	Name  string `json:"name"`
	Alpha string `json:"alpha"`
}

func processStoredLanguages() []languageView {
	loadedLanguages := make([]languageView, len(languages.StoredLanguages))
	i := 0
	for key, lang := range languages.StoredLanguages {
		loadedLanguages[i] = languageView{
			Name:  lang["name"],
			Alpha: key,
		}
		i++
	}

	return loadedLanguages
}

func GetSupportedLanguageHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		if len(languages.StoredLanguages) > 0 {
			return c.JSON(http.StatusOK, processStoredLanguages())
		}

		return c.JSON(http.StatusOK, processStoredLanguages())
	}
}
