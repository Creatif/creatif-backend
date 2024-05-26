package connections

import (
	"creatif/pkg/app/services/publicApi/publicApiError"
	"creatif/pkg/lib/storage"
)

func GetConnections(versionId, projectId, itemId string) ([]Model, error) {
	var models []Model
	if res := storage.Gorm().Raw(getConnectionsSql(), versionId, projectId, itemId, itemId).Scan(&models); res.Error != nil {
		return nil, publicApiError.NewError("getConnections", map[string]string{
			"data": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return models, nil
}

func GetManyConnections(versionId, projectId string, itemIds []string) ([]Model, error) {
	var models []Model
	if res := storage.Gorm().Raw(getConnectionsSql(), versionId, projectId, itemIds, itemIds).Scan(&models); res.Error != nil {
		return nil, publicApiError.NewError("getConnections", map[string]string{
			"data": res.Error.Error(),
		}, publicApiError.DatabaseError)
	}

	return models, nil
}
