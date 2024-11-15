package connections

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/storage"
	"fmt"
)

type QueryVariable struct {
	VariableID    string `gorm:"primarykey;type:text;column:id" json:"variableId"`
	Name          string `json:"name" gorm:"column:name"`
	StructureType string `json:"structureType"`
}

func getConnections(parentVariableId string) ([]declarations.Connection, error) {
	var connections []declarations.Connection
	res := storage.Gorm().Raw(fmt.Sprintf("SELECT * FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), parentVariableId).Scan(&connections)
	if res.Error != nil {
		return nil, res.Error
	}

	return connections, nil
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

func removeParentConnections(variableId string) error {
	res := storage.Gorm().Exec(fmt.Sprintf("DELETE FROM %s WHERE parent_variable_id = ?", (declarations.Connection{}).TableName()), variableId)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
