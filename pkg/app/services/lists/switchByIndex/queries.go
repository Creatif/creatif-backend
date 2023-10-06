package switchByIndex

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
	"gorm.io/gorm"
)

func queryVariableByIndex(projectId, name string, offset int64) (declarations.ListVariable, error) {
	var variable declarations.ListVariable
	res := storage.Gorm().
		Raw(fmt.Sprintf(`
			SELECT lv.id, lv.index
			FROM %s AS lv INNER JOIN %s AS l
			ON l.project_id = ? AND l.name = ? AND lv.list_id = l.id
			OFFSET ? LIMIT 1`, (declarations.ListVariable{}).TableName(), (declarations.List{}).TableName()), projectId, name, offset).
		Scan(&variable)

	if res.Error != nil {
		return declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}

func handleUpdate(source declarations.ListVariable, destination declarations.ListVariable) (declarations.ListVariable, declarations.ListVariable, error) {
	res := storage.Gorm().Exec(fmt.Sprintf(`UPDATE %s SET index = NULL WHERE id = ?`, (declarations.ListVariable{}).TableName()), source.ID)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var toVariable declarations.ListVariable
	res = storage.Gorm().Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? RETURNING id, name, index, short_id`, (declarations.ListVariable{}).TableName()), source.Index, destination.ID).Scan(&toVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var fromVariable declarations.ListVariable
	res = storage.Gorm().Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? RETURNING id, name, index, short_id`, (declarations.ListVariable{}).TableName()), destination.Index, source.ID).Scan(&fromVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return toVariable, fromVariable, nil
}
