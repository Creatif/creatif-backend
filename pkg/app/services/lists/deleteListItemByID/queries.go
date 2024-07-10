package deleteListItemByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getImagePaths(projectId, structureId string) ([]string, error) {
	sql := fmt.Sprintf("SELECT name FROM %s WHERE list_id = ? AND project_id = ?", (declarations.File{}).TableName())

	var images []declarations.File
	if res := storage.Gorm().Raw(sql, structureId, projectId).Scan(&images); res.Error != nil {
		return nil, res.Error
	}

	paths := make([]string, len(images))
	for i, image := range images {
		paths[i] = image.Name
	}

	return paths, nil
}
