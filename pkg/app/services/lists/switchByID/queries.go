package switchByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/app/services/shared"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func queryList(tx *gorm.DB, projectId, name, id, shortID string) (declarations.List, error) {
	id, val := shared.DetermineID("", name, id, shortID)
	var list declarations.List
	if res := tx.Table((declarations.List{}).TableName()).Where(fmt.Sprintf("project_id = ? AND %s", id), projectId, val).Select("ID").First(&list); res.Error != nil {
		return declarations.List{}, res.Error
	}

	return list, nil
}

func queryVariableByID(g *gorm.DB, localeID, listId string, id string, acquireLock bool) (declarations.ListVariable, error) {
	var variable declarations.ListVariable
	lock := ""
	if acquireLock {
		lock = fmt.Sprintf("FOR UPDATE")
	}
	res := g.
		Raw(fmt.Sprintf(`
			SELECT lv.id, lv.index
			FROM %s AS lv WHERE lv.list_id = ? AND id = ? AND locale_id = ? %s`, (declarations.ListVariable{}).TableName(), lock), listId, id, localeID).
		Scan(&variable)

	if res.Error != nil {
		return declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return variable, nil
}

func handleUpdate(g *gorm.DB, source declarations.ListVariable, destination declarations.ListVariable, localeID string) (declarations.ListVariable, declarations.ListVariable, error) {
	res := g.Exec(fmt.Sprintf(`UPDATE %s SET index = NULL WHERE id = ? AND locale_id = ?`, (declarations.ListVariable{}).TableName()), source.ID, localeID)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var toVariable declarations.ListVariable
	res = g.Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? AND locale_id = ? RETURNING id, name, index, short_id, behaviour, groups`, (declarations.ListVariable{}).TableName()), source.Index, destination.ID, localeID).Scan(&toVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	var fromVariable declarations.ListVariable
	res = g.Raw(fmt.Sprintf(`UPDATE %s SET index = ? WHERE id = ? AND locale_id = ? RETURNING id, name, index, short_id, behaviour, groups`, (declarations.ListVariable{}).TableName()), destination.Index, source.ID, localeID).Scan(&fromVariable)
	if res.Error != nil {
		return declarations.ListVariable{}, declarations.ListVariable{}, res.Error
	}

	if res.RowsAffected == 0 {
		return declarations.ListVariable{}, declarations.ListVariable{}, gorm.ErrRecordNotFound
	}

	return toVariable, fromVariable, nil
}

func tryUpdates(projectId, localeID, name, id, shortID, s, d string, currentUpdate, maxUpdates int) (declarations.ListVariable, declarations.ListVariable, error) {
	var to, from declarations.ListVariable
	var lastError error
	if err := storage.Gorm().Transaction(func(tx *gorm.DB) error {
		list, err := queryList(tx, projectId, name, id, shortID)
		if err != nil {
			return err
		}

		source, err := queryVariableByID(tx, localeID, list.ID, s, false)
		if err != nil {
			return err
		}
		destination, err := queryVariableByID(tx, localeID, list.ID, d, false)
		if err != nil {
			return err
		}

		newToVariable, newFromVariable, err := handleUpdate(tx, source, destination, localeID)
		if err != nil {
			return err
		}

		to = newToVariable
		from = newFromVariable

		return nil
	}); err != nil {
		lastError = err
		time.Sleep(time.Millisecond * 10)
		if currentUpdate < maxUpdates {
			return tryUpdates(projectId, localeID, name, id, shortID, s, d, currentUpdate+1, maxUpdates)
		}
	}

	if to.ID == "" || from.ID == "" {
		var err error
		if lastError != nil {
			err = lastError
		} else {
			err = errors.New("Failed switching indexes after retrying.")
		}

		return declarations.ListVariable{}, declarations.ListVariable{}, err
	}
	return to, from, nil
}
