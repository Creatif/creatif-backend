package connections

import (
	"creatif/pkg/app/domain/declarations"
	"encoding/json"
	"github.com/tidwall/sjson"
)

func ReplaceJson(conns []declarations.Connection, value []byte, connectionViewMethod string) ([]byte, error) {
	if connectionViewMethod == "connection" {
		/**
		Get all variables based on the connections
		*/
		connectionVariables, err := getVariableConnectionsView(conns)

		if err != nil {
			return value, err
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
				return nil, err
			}

			updatedValue, err := sjson.SetRawBytes(value, conn.Path, b)
			if err != nil {
				return nil, err
			}

			value = updatedValue
		}
	}

	if connectionViewMethod == "value" {
		/**
		Get all variables based on the connections
		*/
		connectionVariables, err := getVariableValueView(conns)

		if err != nil {
			return value, err
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

			updatedValue, err := sjson.SetRawBytes(value, conn.Path, cv.Value)
			if err != nil {
				return nil, err
			}

			value = updatedValue
		}
	}

	return value, nil
}
