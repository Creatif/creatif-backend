package replaceListItem

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ListWithItem struct {
	ListID    string
	ItemID    string
	ItemIndex time.Time
}

func queryListAndItem(projectId, listName, listId, listShortId, itemID, itemShortID string) (ListWithItem, error) {
	listId, listVal := shared.DetermineID("l", listName, listId, listShortId)
	varId, varVal := shared.DetermineID("lv", "", itemID, itemShortID)
	sql := fmt.Sprintf(`
		SELECT lv.id as item_id, lv.index as item_index, l.id as list_id FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.project_id = ? AND %s AND %s
`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName(), listId, varId)

	var variable ListWithItem
	res := storage.Gorm().Raw(sql, projectId, listVal, varVal).Scan(&variable)
	if res.Error != nil {
		return ListWithItem{}, res.Error
	}

	if res.RowsAffected == 0 {
		return ListWithItem{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}
