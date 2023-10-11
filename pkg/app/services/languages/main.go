package languages

import (
	"creatif/pkg/app/domain/declarations"
	storage2 "creatif/pkg/lib/storage"
)

var StoredLanguages = make(map[string]map[string]string)

func Store() error {
	var languages []declarations.Language
	if res := storage2.Gorm().Find(&languages); res.Error != nil {
		return res.Error
	}

	for _, l := range languages {
		StoredLanguages[l.Alpha] = map[string]string{
			"id":   l.ID,
			"name": l.Name,
		}
	}

	return nil
}
