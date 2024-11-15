package connections

import (
	"encoding/json"
	"github.com/tidwall/sjson"
)

func ReplaceJson(value []byte, variableId string) ([]byte, error) {
	connections, err := getConnections(variableId)
	if err != nil {
		return nil, err
	}

	if len(connections) == 0 {
		return value, nil
	}

	// get all child variables
	for _, c := range connections {
		connectionVariable, err := getChildConnectionVariable(c.ChildStructureType, c.ChildVariableID)
		if err != nil {
			return value, err
		}
		
		b, err := json.Marshal(connectionVariable)
		if err != nil {
			return nil, err
		}

		updatedValue, err := sjson.SetRawBytes(value, c.Path, b)
		if err != nil {
			return nil, err
		}

		return updatedValue, nil
	}

	return value, nil
}
