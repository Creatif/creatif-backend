package appendToList

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

func getHighestIndex(listId string) (float64, error) {
	var idx float64
	res := storage.Gorm().Raw(fmt.Sprintf(`
	SELECT index FROM %s WHERE list_id = ? ORDER BY index DESC LIMIT 1
`, (declarations.ListVariable{}).TableName()),
		listId,
	).Scan(&idx)

	if res.Error != nil {
		return idx, res.Error
	}

	return idx, nil
}
