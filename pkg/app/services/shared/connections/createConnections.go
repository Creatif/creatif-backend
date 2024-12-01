package connections

import (
	"creatif/pkg/app/domain/declarations"
	"creatif/pkg/lib/sdk"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
)

func RemoveJsonConnectionValuePaths(conns []Connection, value []byte) ([]byte, error) {
	for _, c := range conns {
		newValue, err := sjson.DeleteBytes(value, c.Path)
		if err != nil {
			return nil, err
		}

		value = newValue
	}

	return value, nil
}

func CreateConnections(tx *gorm.DB, projectId, parentVariableId, structureType string, conns []Connection, currentValue []byte) ([]byte, []declarations.Connection, error) {
	newValue, err := RemoveJsonConnectionValuePaths(conns, currentValue)
	if err != nil {
		return nil, nil, err
	}

	if err := CheckConnectionsIntegrity(tx, conns); err != nil {
		return nil, nil, err
	}

	created := sdk.Map(conns, func(idx int, value Connection) declarations.Connection {
		return declarations.NewConnection(
			projectId,
			value.Path,
			parentVariableId,
			structureType,
			value.VariableID,
			value.StructureType,
		)
	})

	return newValue, created, nil
}
