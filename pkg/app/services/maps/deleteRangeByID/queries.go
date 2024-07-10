package deleteRangeByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getImagePaths(projectId string, structureIds []string) ([]string, error) {
	sql := fmt.Sprintf("SELECT name FROM %s WHERE map_id IN(?) AND project_id = ?", (declarations.File{}).TableName())

	var images []declarations.File
	if res := storage.Gorm().Raw(sql, structureIds, projectId).Scan(&images); res.Error != nil {
		return nil, res.Error
	}

	paths := make([]string, len(images))
	for i, image := range images {
		paths[i] = image.Name
	}

	return paths, nil
}
