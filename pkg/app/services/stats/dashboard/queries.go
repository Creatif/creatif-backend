package dashboard

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getLists() ([]LogicModel, error) {
	sql := fmt.Sprintf("SELECT COUNT(id) as count, name, created_at, 'list' AS type FROM %s GROUP BY id", (declarations.List{}).TableName())

	var models []LogicModel
	if res := storage.Gorm().Raw(sql).Scan(&models); res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}

func getMaps() ([]LogicModel, error) {
	sql := fmt.Sprintf("SELECT COUNT(id) as count, name, created_at, 'map' AS type FROM %s GROUP BY id", (declarations.Map{}).TableName())

	var models []LogicModel
	if res := storage.Gorm().Raw(sql).Scan(&models); res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}
