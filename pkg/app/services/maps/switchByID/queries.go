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
	source      declarations.MapVariable
	destination declarations.MapVariable
}

/*
*
Gets the highest and lowest index in the entire variable list.
*/
func getHighestLowestIndex(mapId, projectId string) (indexRange, error) {
	model := indexRange{}
	res := storage.Gorm().Raw(fmt.Sprintf(`
	SELECT (SELECT index FROM %s WHERE map_id = ? ORDER BY index DESC LIMIT 1) AS highest,
		(SELECT index FROM %s WHERE map_id = ? ORDER BY index ASC LIMIT 1) AS lowest
`, (declarations.MapVariable{}).TableName(), (declarations.MapVariable{}).TableName()),
		mapId, mapId,
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
	var s declarations.MapVariable
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT lv.id, lv.index FROM %s AS lv INNER JOIN %s AS l ON lv.map_id = l.id AND l.project_id = ? AND (l.id = ? OR l.name = ? OR l.short_id = ?) AND lv.id = ?",
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
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

	var d declarations.MapVariable
	res = storage.Gorm().Raw(fmt.Sprintf("SELECT lv.id, lv.index FROM %s AS lv INNER JOIN %s AS l ON lv.map_id = l.id AND l.project_id = ? AND (l.id = ? OR l.name = ? OR l.short_id = ?) AND lv.id = ?",
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
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

func updateWithCustomIndex(idx float64, id, mapId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf(`
UPDATE %s
SET index = ? WHERE id = ? AND map_id = ?
`,
		(declarations.MapVariable{}).TableName(),
	), idx, id, mapId)

	if res.Error != nil {
		return appErrors.NewApplicationError(res.Error)
	}

	if res.RowsAffected == 0 {
		return appErrors.NewNotFoundError(errors.New("Could not switch map variables."))
	}

	return nil
}

/*
*
Gets the index the one before the destination.
*/
func getIndexDesc(mapId string, destinationIndex float64) (float64, error) {
	sql := fmt.Sprintf("SELECT index FROM %s AS vg WHERE vg.map_id = ? AND vg.index > ? ORDER BY vg.index ASC LIMIT 1", (declarations.MapVariable{}).TableName())
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
    WHERE vg.map_id = ?
      AND vg.index < ?
    ORDER BY vg.index ASC
) AS sorted_results
ORDER BY index DESC 
LIMIT 1;
`, (declarations.MapVariable{}).TableName())
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
	sql := fmt.Sprintf("UPDATE %s SET index = ? WHERE id = ? AND map_id = ?", (declarations.MapVariable{}).TableName())

	if res := storage.Gorm().Exec(sql, index, destinationVariableId, mapId); res.Error != nil {
		return res.Error
	}

	return nil
}
