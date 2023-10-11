package declarations

import (
	app2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	storage2 "creatif/pkg/lib/storage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type languageView struct {
	Name  string `json:"name"`
	Alpha string `json:"alpha"`
}

var loadedLanguages []languageView

func GetSupportedLanguageHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		if len(loadedLanguages) > 0 {
			return c.JSON(http.StatusOK, loadedLanguages)
		}

		var languages []app2.Language
		res := storage2.Gorm().Find(&languages)
		if res.Error != nil {
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}

		if res.RowsAffected == 0 {
			return c.JSON(http.StatusInternalServerError, "Internal server error")
		}

		loadedLanguages = sdk.Map(languages, func(idx int, value app2.Language) languageView {
			return languageView{
				Name:  value.Name,
				Alpha: value.Alpha,
			}
		})

		return c.JSON(http.StatusOK, loadedLanguages)
	}
}
