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
	VariableID             string `json:"variableId"`
	Value                  string `json:"value"`
	Path                   string `json:"path"`
	StructureType          string `json:"structureType"`
	CreatifSpecialVariable bool   `json:"creatif_special_variable"`
}

func getBulkVariablesFromConnections(conns []declarations.Connection) ([]QueryVariable, error) {
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

	var bufferVariables []QueryVariable
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT id, name, 'map' AS structure_type FROM %s WHERE id IN(?)", (declarations.MapVariable{}).TableName()), mapVariableIds).Scan(&bufferVariables)
	if res.Error != nil {
		return nil, res.Error
	}

	var listVariables []QueryVariable
	res = storage.Gorm().Raw(fmt.Sprintf("SELECT id, name, 'list' AS structure_type FROM %s WHERE id IN(?)", (declarations.ListVariable{}).TableName()), listVariableIds).Scan(&listVariables)
	if res.Error != nil {
		return nil, res.Error
	}

	bufferVariables = append(bufferVariables, listVariables...)

	return bufferVariables, nil
}

func RemoveParentConnections(variableId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), variableId)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
