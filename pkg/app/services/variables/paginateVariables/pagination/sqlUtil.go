package pagination

import (
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

func getInitialID(projectId, table, orderBy string) (string, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE project_id = ? ORDER BY id %s LIMIT 1", table, orderBy)
	var model initialModel
	if res := storage.Gorm().Raw(sql, projectId).Scan(&model); res.Error != nil {
		return "", res.Error
	}

	return model.ID, nil
}

func getOperator(direction, orderBy string, isFirst bool) (string, error) {
	if direction == DIRECTION_FORWARD && isFirst && orderBy == DESC {
		return "<=", nil
	} else if direction == DIRECTION_FORWARD && !isFirst && orderBy == DESC {
		return "<", nil
	} else if direction == DIRECTION_FORWARD && isFirst && orderBy == ASC {
		return ">=", nil
	} else if direction == DIRECTION_FORWARD && !isFirst && orderBy == ASC {
		return ">", nil
	} else if direction == DIRECTION_BACKWARDS && isFirst && orderBy == ASC {
		return ">=", nil
	} else if direction == DIRECTION_BACKWARDS && !isFirst && orderBy == ASC {
		return ">", nil
	}

	return "", errors.New("Operator could not be determined")
}
