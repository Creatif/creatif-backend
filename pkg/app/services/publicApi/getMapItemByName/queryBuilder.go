package getMapItemByName

import (
	"creatif/pkg/app/services/locales"
)

func createPlaceholders(projectId, versionId, structureName, variableName, locale string) (map[string]interface{}, error) {
	placeholders := make(map[string]interface{})
	placeholders["projectId"] = projectId
	placeholders["versionId"] = versionId
	placeholders["structureName"] = structureName
	placeholders["variableName"] = variableName

	if locale == "" {
		l, _ := locales.GetIDWithAlpha("eng")
		locale = l
		placeholders["localeId"] = l
	}

	if locale != "" {
		l, _ := locales.GetIDWithAlpha(locale)
		locale = l
		placeholders["localeId"] = l
	}

	return placeholders, nil
}
