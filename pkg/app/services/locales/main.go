package locales

import (
	"creatif/pkg/app/domain/declarations"
	storage2 "creatif/pkg/lib/storage"
	"errors"
)

var storedLocales = make(map[string]map[string]string)
var idToAlpha = make(map[string]string)

var LocaleByAlphaNotExists = errors.New("Locale does not exist")
var IDByAlphaNotExists = errors.New("ID by alpha does not exist")

func Store() error {
	if len(storedLocales) > 0 {
		return nil
	}

	var locales []declarations.Locale
	if res := storage2.Gorm().Find(&locales); res.Error != nil {
		return res.Error
	}

	for _, l := range locales {
		storedLocales[l.Alpha] = map[string]string{
			"id":   l.ID,
			"name": l.Name,
		}

		idToAlpha[l.ID] = l.Alpha
	}

	return nil
}

func ExistsByAlpha(alpha string) bool {
	_, ok := storedLocales[alpha]

	return ok
}

func Locales() map[string]map[string]string {
	return storedLocales
}

func GetIDWithAlpha(alpha string) (string, error) {
	if !ExistsByAlpha(alpha) {
		return "", LocaleByAlphaNotExists
	}

	val, _ := storedLocales[alpha]
	return val["id"], nil
}

func GetAlphaWithID(id string) (string, error) {
	val, ok := idToAlpha[id]
	if !ok {
		return "", IDByAlphaNotExists
	}

	return val, nil
}
