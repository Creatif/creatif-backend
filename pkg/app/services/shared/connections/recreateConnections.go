package connections

import (
	"creatif/pkg/app/domain/declarations"
	"gorm.io/gorm"
)

func RecreateConnections(tx *gorm.DB, projectId, parentVariableId, structureType string, conns []Connection, currentValue []byte) ([]byte, []declarations.Connection, error) {
	if err := RemoveParentConnections(tx, parentVariableId); err != nil {
		return nil, nil, err
	}

	return CreateConnections(tx, projectId, parentVariableId, structureType, conns, currentValue)
}
