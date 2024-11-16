package connections

import (
	"creatif/pkg/app/domain/declarations"
	"encoding/json"
	"github.com/tidwall/sjson"
)

func ReplaceJson(value []byte, variableId string) ([]declarations.Connection, []byte, error) {
	connections, err := getConnections(variableId)
	if err != nil {
		return nil, nil, err
	}

	if len(connections) == 0 {
		return nil, value, nil
	}

	// get all child variables
	for _, c := range connections {
		connectionVariable, err := getChildConnectionVariable(c.ChildStructureType, c.ChildVariableID)
		if err != nil {
			return nil, value, err
		}

		b, err := json.Marshal(connectionVariable)
		if err != nil {
			return nil, nil, err
		}

		updatedValue, err := sjson.SetRawBytes(value, c.Path, b)
		if err != nil {
			return nil, nil, err
		}

		return nil, updatedValue, nil
	}

	return connections, value, nil
}
