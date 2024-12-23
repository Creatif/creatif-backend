package switchByID

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
)

type indexRange struct {
	Highest float64
	Lowest  float64
}

type sourceDestinationVariables struct {
	source      declarations.ListVariable
	destination declarations.ListVariable
}

func getHighestLowestIndex(listId, projectId string) (indexRange, error) {
	model := indexRange{}
	res := storage.Gorm().Raw(fmt.Sprintf(`
	SELECT (SELECT index FROM %s WHERE list_id = ? ORDER BY index DESC LIMIT 1) AS highest,
		(SELECT index FROM %s WHERE list_id = ? ORDER BY index ASC LIMIT 1) AS lowest FROM
	%s WHERE list_id = ?
`, (declarations.ListVariable{}).TableName(), (declarations.ListVariable{}).TableName(), (declarations.ListVariable{}).TableName()),
		listId, listId, listId,
	).Scan(&model)

	if res.Error != nil {
		return indexRange{}, res.Error
	}

	if res.RowsAffected == 0 {
		return indexRange{}, errors.New("not_found")
	}

	return model, nil
}

func getSourceDestinationVariables(projectId, name, source, destination string) (sourceDestinationVariables, error) {
	var s declarations.ListVariable
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT lv.id, lv.index FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.project_id = ? AND (l.id = ? OR l.name = ? OR l.short_id = ?) AND lv.id = ?",
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
	), projectId, name, name, name, source).Scan(&s)

	if res.Error != nil {
		return sourceDestinationVariables{}, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Source or destination do not exist",
		})
	}

	if res.RowsAffected == 0 {
		return sourceDestinationVariables{}, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Source or destination do not exist",
		})
	}

	var d declarations.ListVariable
	res = storage.Gorm().Raw(fmt.Sprintf("SELECT lv.id, lv.index FROM %s AS lv INNER JOIN %s AS l ON lv.list_id = l.id AND l.project_id = ? AND (l.id = ? OR l.name = ? OR l.short_id = ?) AND lv.id = ?",
		(declarations.ListVariable{}).TableName(),
		(declarations.List{}).TableName(),
	), projectId, name, name, name, destination).Scan(&d)

	if res.Error != nil {
		return sourceDestinationVariables{}, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Source or destination do not exist",
		})
	}

	if res.RowsAffected == 0 {
		return sourceDestinationVariables{}, appErrors.NewValidationError(map[string]string{
			"invalidSourceDestination": "Source or destination do not exist",
		})
	}

	return sourceDestinationVariables{source: s, destination: d}, nil
}

func updateWithCustomIndex(idx float64, id, listId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf(`
UPDATE %s
SET index = ? WHERE id = ? AND list_id = ?
`,
		(declarations.ListVariable{}).TableName(),
	), idx, id, listId)

	if res.Error != nil {
		return appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return appErrors.NewNotFoundError(errors.New("Could not switch list variables."))
	}

	return nil
}

/*
*
Gets the index the one before the destination.
*/
func getIndexDesc(mapId string, destinationIndex float64) (float64, error) {
	sql := fmt.Sprintf("SELECT index FROM %s AS vg WHERE vg.list_id = ? AND vg.index > ? ORDER BY vg.index ASC LIMIT 1", (declarations.ListVariable{}).TableName())
	var idx float64
	res := storage.Gorm().Raw(sql, mapId, destinationIndex).Scan(&idx)

	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, errors.New("Destination upper index not found")
	}

	return idx, nil
}

func getIndexAsc(mapId string, destinationIndex float64) (float64, error) {
	sql := fmt.Sprintf(`
SELECT index 
FROM (
    SELECT index 
    FROM %s AS vg 
    WHERE vg.list_id = ?
      AND vg.index < ?
    ORDER BY vg.index ASC
) AS sorted_results
ORDER BY index DESC 
LIMIT 1;
`, (declarations.ListVariable{}).TableName())
	var idx float64
	res := storage.Gorm().Raw(sql, mapId, destinationIndex).Scan(&idx)

	if res.Error != nil {
		return 0, res.Error
	}

	if res.RowsAffected == 0 {
		return 0, errors.New("Destination upper index not found")
	}

	return idx, nil
}

func updateDestinationIndex(mapId, destinationVariableId string, index float64) error {
	sql := fmt.Sprintf("UPDATE %s SET index = ? WHERE id = ? AND list_id = ?", (declarations.ListVariable{}).TableName())

	if res := storage.Gorm().Exec(sql, index, destinationVariableId, mapId); res.Error != nil {
		return res.Error
	}

	return nil
}
