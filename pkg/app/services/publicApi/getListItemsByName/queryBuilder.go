package getListItemsByName

import (
	"creatif/pkg/app/services/locales"
	"creatif/pkg/app/services/publicApi/publicApiError"
)

func createPlaceholders(projectId, versionId, structureName, variableName, locale string) (map[string]interface{}, error) {
	placeholders := make(map[string]interface{})
	placeholders["projectId"] = projectId
	placeholders["versionId"] = versionId
	placeholders["structureName"] = structureName
	placeholders["variableName"] = variableName

	if locale == "" {
		l, _ := locales.GetIDWithAlpha("eng")
		placeholders["localeId"] = l
	}

	if locale != "" {
		l, err := locales.GetIDWithAlpha(locale)
		if err != nil {
			return nil, publicApiError.NewError("getListItemsByName", map[string]string{
				"invalidLocale": "The locale you provided is invalid and does not exist.",
			}, publicApiError.ValidationError)
		}

		placeholders["localeId"] = l
	}

	return placeholders, nil
}
