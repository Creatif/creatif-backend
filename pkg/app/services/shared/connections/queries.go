package connections

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

type QueryVariable struct {
	VariableID    string `gorm:"primarykey;type:text;column:id" json:"variableId"`
	Name          string `json:"name" gorm:"column:name"`
	StructureType string `json:"structureType" gorm:"column:structure_type"`
}

type ConnectionVariable struct {
	VariableID    string `json:"variableId"`
	Value         string `json:"value"`
	Path          string `json:"path"`
	StructureType string `json:"structureType"`
}

func getChildConnectionFromParent(parentVariableId string) ([]declarations.Connection, error) {
	var connections []declarations.Connection
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT * FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), parentVariableId).Scan(&connections)
	if res.Error != nil {
		return nil, res.Error
	}

	return connections, nil
}

func getBulkConnectionVariablesFromConnections(conns []declarations.Connection) ([]QueryVariable, error) {
	mapVariableIds := make([]string, 0)
	listVariableIds := make([]string, 0)

	for _, c := range conns {
		if c.ChildStructureType == "map" {
			mapVariableIds = append(mapVariableIds, c.ChildVariableID)
		}

		if c.ChildStructureType == "list" {
			listVariableIds = append(listVariableIds, c.ChildVariableID)
		}
	}

	var mapVariables []QueryVariable
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name, 'map' AS structure_type FROM %s WHERE id IN(?)", (declarations.MapVariable{}).TableName()), mapVariableIds).Scan(&mapVariables)
	if res.Error != nil {
		return nil, res.Error
	}

	var listVariables []QueryVariable
	res = storage.Gorm().Raw(fmt.Sprintf("SELECT id, name, 'list' AS structure_type FROM %s WHERE id IN(?)", (declarations.ListVariable{}).TableName()), listVariableIds).Scan(&listVariables)
	if res.Error != nil {
		return nil, res.Error
	}

	mapVariables = append(mapVariables, listVariables...)

	return mapVariables, nil
}

func getChildConnectionVariable(childStructureType, childVariableId string) (QueryVariable, error) {
	var connectionVariable QueryVariable
	if childStructureType == "map" {
		res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name from %s WHERE id = ?", (declarations.MapVariable{}).TableName()), childVariableId).Scan(&connectionVariable)
		if res.Error != nil {
			return QueryVariable{}, res.Error
		}

		connectionVariable.StructureType = "map"
	}

	if childStructureType == "list" {
		res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name from %s WHERE id = ?", (declarations.ListVariable{}).TableName()), childVariableId).Scan(&connectionVariable)
		if res.Error != nil {
			return QueryVariable{}, res.Error
		}

		connectionVariable.StructureType = "list"
	}

	return connectionVariable, nil
}

func RemoveParentConnections(variableId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), variableId)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
