package pagination

import (
	"creatif/pkg/lib/storage"
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

func getInitialOperator(direction, orderByDirection string) string {
	if direction == DIRECTION_FORWARD && orderByDirection == DESC {
		return "<="
	} else if direction == DIRECTION_FORWARD && orderByDirection == ASC {
		return ">="
	} else if direction == DIRECTION_BACKWARDS && orderByDirection == DESC {
		return ""
	}

	return ""
}

func getOperator(direction, orderBy string) string {
	if direction == DIRECTION_FORWARD && orderBy == DESC {
		return "<"
	} else if direction == DIRECTION_FORWARD && orderBy == ASC {
		return ">"
	} else if direction == DIRECTION_BACKWARDS && orderBy == DESC {
		return ""
	}

	return ""
}
