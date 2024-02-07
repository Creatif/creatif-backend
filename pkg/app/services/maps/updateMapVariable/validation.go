package updateMapVariable

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/appErrors"
	"creatif/pkg/lib/constants"
	"creatif/pkg/lib/logger"
	"creatif/pkg/lib/storage"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func validateGroupsNumAndBehaviour(mapName, projectId, variableName string, groups []string, logBuilder logger.LogBuilder) error {
	type GroupBehaviourCheck struct {
		Count     int    `gorm:"column:count"`
		Behaviour string `gorm:"column:behaviour"`
	}

	var check GroupBehaviourCheck
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT cardinality(mv.groups) AS count, behaviour
FROM %s AS mv 
INNER JOIN %s AS m ON (m.name = ? OR m.id = ? OR m.short_id = ?) AND m.project_id = ? AND m.id = mv.map_id AND (mv.id = ? OR mv.short_id = ?)`,
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName()),
		mapName,
		mapName,
		mapName,
		projectId,
		variableName,
		variableName,
	).Scan(&check)

	if res.Error != nil || res.RowsAffected == 0 {
		if res.Error != nil {
			logBuilder.Add("updateMapVariable", res.Error.Error())
		} else {
			logBuilder.Add("updateMapVariable", "No rows returned. Might be a bug")
		}
		return appErrors.NewValidationError(map[string]string{
			"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", variableName),
		})
	}

	if len(groups) > 0 {
		if check.Count+len(groups) > 20 {
			return appErrors.NewValidationError(map[string]string{
				"groups": fmt.Sprintf("Invalid number of groups for '%s'. Maximum number of groups per variable is 20.", variableName),
			})
		}
	}

	if check.Behaviour == constants.ReadonlyBehaviour {
		return appErrors.NewValidationError(map[string]string{
			"behaviourReadonly": fmt.Sprintf("Map item with ID '%s' is readonly and cannot be updated.", variableName),
		})
	}

	return nil
}

func validateUniqueName(mapId, variableId, mapVariableName, projectId string) error {
	var id string
	res := storage.Gorm().Raw(fmt.Sprintf(`
SELECT mv.id
FROM %s AS mv 
INNER JOIN %s AS m ON (m.name = ? OR m.id = ? OR m.short_id = ?) AND m.project_id = ? AND m.id = mv.map_id AND mv.name = ? AND mv.id != ?`,
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName()),
		mapId,
		mapId,
		mapId,
		projectId,
		mapVariableName,
		variableId,
	).Scan(&id)

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable map item with name '%s' already exists.", mapVariableName),
		})
	}

	if res.RowsAffected != 0 {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable map item with name '%s' already exists.", mapVariableName),
		})
	}

	if id != "" {
		return appErrors.NewValidationError(map[string]string{
			"exists": fmt.Sprintf("Variable map item with name '%s' already exists.", mapVariableName),
		})
	}

	return nil
}
