package addToMap

import (
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

type parentReference struct {
	ParentID string `gorm:"column:list_id"`
	ID       string `gorm:"column:id"`
}

func getParentReferences(structureName, structureType, variableId, parentId string) (parentReference, error) {
	if structureType == "list" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS parent_id FROM declarations.list_variables AS lv INNER JOIN declarations.lists AS l ON l.id = lv.list_id AND l.name = ? AND lv.id = ?`)

		var pr parentReference
		if res := storage.Gorm().Raw(sql, structureName, variableId).Scan(&pr); res.Error != nil {
			return parentReference{}, res.Error
		}

		if pr.ParentID == parentId {
			return parentReference{}, errors.New("Invalid parent reference. A reference can not have itself as the parent")
		}

		return pr, nil
	}

	if structureType == "map" {
		sql := fmt.Sprintf(`SELECT lv.id AS id, l.id AS parent_id FROM declarations.map_variables AS lv INNER JOIN declarations.maps AS l ON l.id = lv.map_id AND l.name = ? AND lv.id = ?`)

		var pr parentReference
		if res := storage.Gorm().Raw(sql, structureName, variableId).Scan(&pr); res.Error != nil {
			return parentReference{}, res.Error
		}

		if pr.ParentID == parentId {
			return parentReference{}, errors.New("Invalid parent reference. A reference can not have itself as the parent")
		}

		return pr, nil
	}

	sql := fmt.Sprintf(`SELECT id FROM declarations.variables WHERE id = ?`)

	var pr parentReference
	if res := storage.Gorm().Raw(sql, variableId).Scan(&pr); res.Error != nil {
		return pr, res.Error
	}

	return pr, nil
}
