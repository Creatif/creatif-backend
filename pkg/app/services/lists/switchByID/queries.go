package switchByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func queryList(tx *gorm.DB, projectId, name string) (declarations.List, error) {
	var list declarations.List
	if res := tx.Table((declarations.List{}).TableName()).Where("project_id = ? AND name = ?", projectId, name).Select("ID").First(&list); res.Error != nil {
		return declarations.List{}, res.Error
	}

	return list, nil
}

func queryVariableByID(g *gorm.DB, listId string, id string, acquireLock bool) (declarations.ListVariable, error) {
	var variable declarations.ListVariable
	lock := ""
	if acquireLock {
		lock = fmt.Sprintf("FOR UPDATE")
	}
	res := g.
		Raw(fmt.Sprintf(`
			SELECT lv.id, lv.index
			FROM %s AS lv WHERE lv.list_id = ? AND id = ? %s`, (declarations.ListVariable{}).TableName(), lock), listId, id).
		Scan(&variable)

	if res.Error != nil {
		return declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}

func handleUpdate(g *gorm.DB, source declarations.ListVariable, destination declarations.ListVariable) (declarations.ListVariable, declarations.ListVariable, error) {
	res := g.Exec(fmt.Sprintf(`UPDATE %s SET index = NULL WHERE id = ?`, (declarations.ListVariable{}).TableName()), source.ID)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var toVariable declarations.ListVariable
	res = g.Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? RETURNING id, name, index, short_id`, (declarations.ListVariable{}).TableName()), source.Index, destination.ID).Scan(&toVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var fromVariable declarations.ListVariable
	res = g.Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? RETURNING id, name, index, short_id`, (declarations.ListVariable{}).TableName()), destination.Index, source.ID).Scan(&fromVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return toVariable, fromVariable, nil
}

func tryUpdates(projectId, name, s, d string, currentUpdate, maxUpdates int) (declarations.ListVariable, declarations.ListVariable, error) {
	var to, from declarations.ListVariable
	if err := storage.Gorm().Transaction(func(tx *gorm.DB) error {
		list, err := queryList(tx, projectId, name)
		if err != nil {
			return err
		}

		source, err := queryVariableByID(tx, list.ID, s, false)
		if err != nil {
			return err
		}
		destination, err := queryVariableByID(tx, list.ID, d, false)
		if err != nil {
			return err
		}

		newToVariable, newFromVariable, err := handleUpdate(tx, source, destination)
		if err != nil {
			return err
		}

		to = newToVariable
		from = newFromVariable

		return nil
	}); err != nil {
		time.Sleep(time.Millisecond * 10)
		if currentUpdate < maxUpdates {
			return tryUpdates(projectId, name, s, d, currentUpdate+1, maxUpdates)
		}
	}

	if to.ID == "" || from.ID == "" {
		return declarations.ListVariable{}, declarations.ListVariable{}, errors.New("Failed switching indexes after retrying.")
	}
	return to, from, nil
}
