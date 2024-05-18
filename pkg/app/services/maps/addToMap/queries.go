package addToMap

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getHighestIndex(mapId string) (float64, error) {
	var highestIndex float64
	res := storage.Gorm().Raw(fmt.Sprintf(`
	SELECT COALESCE(MAX(index), 0) AS index FROM %s WHERE map_id = ?
`, (declarations.MapVariable{}).TableName()),
		mapId,
	).Scan(&highestIndex)

	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, nil
	}

	return highestIndex, nil
}
