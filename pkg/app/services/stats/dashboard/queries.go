package dashboard

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/domain/published"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getLists(projectId string) ([]StructureLogicModel, error) {
	sql := fmt.Sprintf("SELECT v.id, (SELECT COUNT(vr.id) AS count FROM %s AS vr WHERE vr.list_id = v.id) AS count, v.name, v.created_at, 'list' AS type FROM %s AS v WHERE v.project_id = ? GROUP BY v.id ORDER BY v.created_at DESC", (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName())

	var models []StructureLogicModel
	if res := storage.Gorm().Raw(sql, projectId).Scan(&models); res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}

func getMaps(projectId string) ([]StructureLogicModel, error) {
	sql := fmt.Sprintf("SELECT v.id, (SELECT COUNT(vr.id) AS count FROM %s AS vr WHERE vr.map_id = v.id) AS count, v.name, v.created_at, 'map' AS type FROM %s AS v WHERE v.project_id = ? GROUP BY v.id ORDER BY v.created_at DESC", (declarations.MapVariable{}).TableName(), (declarations.Map{}).TableName())

	var models []StructureLogicModel
	if res := storage.Gorm().Raw(sql, projectId).Scan(&models); res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}

func getVersions(projectId string) ([]VersionLogicModel, error) {
	sql := fmt.Sprintf("SELECT id, name, project_id, created_at, updated_at, is_production_version FROM %s WHERE project_id = ? ORDER BY created_at DESC", (published.Version{}).TableName())

	var versions []VersionLogicModel
	if res := storage.Gorm().Raw(sql, projectId).Scan(&versions); res.Error != nil {
		return nil, res.Error
	}

	return versions, nil
}
