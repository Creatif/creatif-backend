package app

import (
	app2 "creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	storage2 "creatif/pkg/lib/storage"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type languageView struct {
	Name    string `json:"name"`
	Alpha3b string `json:"alpha3b"`
	Alpha3t string `json:"alpha3t"`
	Alpha2  string `json:"alpha2"`
}

var loadedLanguages []languageView

func GetSupportedLanguageHandler() func(e echo.Context) error {
	return func(c echo.Context) error {
		if len(loadedLanguages) > 0 {
			fmt.Println(len(loadedLanguages))
			fmt.Println("Already loaded")
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

		fmt.Println("Not loaded")
		loadedLanguages = sdk.Map(languages, func(idx int, value app2.Language) languageView {
			return languageView{
				Name:    value.Name,
				Alpha3b: value.Alpha3b,
				Alpha3t: value.Alpha3t,
				Alpha2:  value.Alpha2,
			}
		})

		return c.JSON(http.StatusOK, loadedLanguages)
	}
}
