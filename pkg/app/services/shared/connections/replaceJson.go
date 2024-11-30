package connections

import (
	"creatif/pkg/app/domain/declarations"
	"encoding/json"
	"github.com/tidwall/sjson"
)

func ReplaceJson(value []byte, variableId string) ([]declarations.Connection, []byte, error) {
	/**
	Gets all connections where variableId is the parent. This will effectively return
	all child connections of this variableId.
	*/
	conns, err := getChildConnectionFromParent(variableId)

	if err != nil {
		return nil, nil, err
	}

	if len(conns) == 0 {
		return nil, value, nil
	}

	/**
	Get all variables based on the connections
	*/
	connectionVariables, err := getBulkConnectionVariablesFromConnections(conns)

	if err != nil {
		return nil, value, err
	}

	/**
	Connections and its variables are linked. We need data from both to construct ConnectionVariable struct.
	Because of that, we can group them into a single structure (a map) for faster lookup in order to not
	do that in the below iteration.
	*/
	connsMap := make(map[string]declarations.Connection)
	for _, c := range conns {
		connsMap[c.ChildVariableID] = c
	}

	for _, cv := range connectionVariables {
		conn := connsMap[cv.VariableID]

		viewConnection := ConnectionVariable{
			VariableID:             cv.VariableID,
			Value:                  cv.Name,
			Path:                   conn.Path,
			StructureType:          cv.StructureType,
			CreatifSpecialVariable: true,
		}

		b, err := json.Marshal(viewConnection)
		if err != nil {
			return nil, nil, err
		}

		updatedValue, err := sjson.SetRawBytes(value, conn.Path, b)
		if err != nil {
			return nil, nil, err
		}

		value = updatedValue
	}

	return conns, value, nil
}
