package paginateReferences

import (
	"creatif/pkg/app/domain/declarations"
	"fmt"
	"strings"
)

func getWorkingTables(structureType string) [2]string {
	tables := [2]string{
		(declarations.MapVariable{}).TableName(),
		(declarations.Map{}).TableName(),
	}

	if structureType == "list" {
		tables = [2]string{
			(declarations.List{}).TableName(),
			(declarations.ListVariable{}).TableName(),
		}
	}

	return tables
}

func createPlaceholdersFromModel(model Model) map[string]interface{} {
	placeholders := make(map[string]interface{})
	placeholders["projectID"] = model.ProjectID
	placeholders["offset"] = (model.Page - 1) * model.Limit
	placeholders["limit"] = model.Limit
	placeholders["parentReference"] = model.ParentID
	placeholders["childReference"] = model.ChildID
	placeholders["childStructureID"] = model.ChildStructureID
	placeholders["parentStructureID"] = model.ParentStructureID

	if model.Behaviour != "" {
		placeholders["behaviour"] = model.Behaviour
	}

	if model.Search != "" {
		placeholders["searchOne"] = fmt.Sprintf("%%%s", model.Search)
		placeholders["searchTwo"] = fmt.Sprintf("%s%%", model.Search)
		placeholders["searchThree"] = fmt.Sprintf("%%%s%%", model.Search)
		placeholders["searchFour"] = model.Search
	}

	return placeholders
}

func createFields(model Model) (string, string) {
	orderBy, direction := model.OrderBy, model.OrderDirection
	if model.OrderBy == "" {
		orderBy = "index"
	}

	if model.OrderDirection == "" {
		direction = "ASC"
	}

	direction = strings.ToUpper(direction)
	return orderBy, direction
}
