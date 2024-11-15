package connections

import "creatif/pkg/app/domain/declarations"

func RecreateConnections(projectId, parentVariableId, structureType string, conns []Connection, currentValue []byte) ([]byte, []declarations.Connection, error) {
	if err := removeParentConnections(parentVariableId); err != nil {
		return nil, nil, err
	}

	return CreateConnections(projectId, parentVariableId, structureType, conns, currentValue)
}
