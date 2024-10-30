package create

import (
	"creatif/pkg/app/domain/app"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/datatypes"
)

type DataQuery struct {
	ID   string         `gorm:"column:id"`
	Data datatypes.JSON `gorm:"column:data"`
}

func getActivityCount(projectId string) (int64, error) {
	sql := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE project_id = ?", (app.Activity{}).TableName())

	var count int64
	if res := storage.Gorm().Raw(sql, projectId).Scan(&count); res.Error != nil {
		return 0, res.Error
	}

	return count, nil
}

func getActivities(projectId string) ([]DataQuery, error) {
	sql := fmt.Sprintf("SELECT id, data FROM %s WHERE project_id = ? ORDER BY created_at DESC", (app.Activity{}).TableName())

	var dataQuery []DataQuery
	if res := storage.Gorm().Raw(sql, projectId).Scan(&dataQuery); res.Error != nil {
		return nil, res.Error
	}

	return dataQuery, nil
}

func getLastActivityDataQuery(projectId string) (DataQuery, error) {
	sql := fmt.Sprintf("SELECT id, data FROM %s WHERE project_id = ? ORDER BY created_at DESC LIMIT 1", (app.Activity{}).TableName())

	var dataQuery DataQuery
	if res := storage.Gorm().Raw(sql, projectId).Scan(&dataQuery); res.Error != nil {
		return DataQuery{}, res.Error
	}

	return dataQuery, nil
}

func deleteActivity(projectId, id string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE project_id = ? AND id = ?", (app.Activity{}).TableName())

	var count int64
	if res := storage.Gorm().Raw(sql, projectId, id).Scan(&count); res.Error != nil {
		return res.Error
	}

	return nil
}
