package services

import (
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
)

func Execute(mapId string, strategy QueryStrategy) ([]map[string]interface{}, error) {
	models := make([]map[string]interface{}, 0)
	if res := storage.Gorm().Raw(strategy.GetQuery(), mapId).Scan(&models); res.Error != nil {
		return nil, appErrors.NewDatabaseError(res.Error).AddError("getMap.Services.Execute", nil)
	}

	return models, nil
}
