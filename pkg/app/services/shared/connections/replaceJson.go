package connections

import (
	"creatif/pkg/app/domain/declarations"
	"encoding/json"
	"github.com/tidwall/sjson"
)

func ReplaceJson(value []byte, variableId string) ([]declarations.Connection, []byte, error) {
	conns, err := getConnections(variableId)
	if err != nil {
		return nil, nil, err
	}

	if len(conns) == 0 {
		return nil, value, nil
	}

	// get all child variables
	for _, c := range conns {
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

		value = updatedValue
	}

	return conns, value, nil
}
