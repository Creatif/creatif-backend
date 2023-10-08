package replaceListItem

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
)

type ListWithItem struct {
	ListID    string
	ItemID    string
	ItemIndex string
}

func queryListAndItem(projectId, listName, itemName string) (ListWithItem, error) {
	sql := fmt.Sprintf(`
		SELECT lv.id as item_id, lv.index as item_index, l.id as list_id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.project_id = ? AND l.name = ? AND lv.name = ?
`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName())

	var variable ListWithItem
	res := storage.Gorm().Raw(sql, projectId, listName, itemName).Scan(&variable)
	if res.Error != nil {
		return ListWithItem{}, res.Error
	}

	if res.RowsAffected == 0 {
		return ListWithItem{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}
