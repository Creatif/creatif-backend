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

		for _, c := range conns {
			var connectionVariable QueryConnectionView
			for _, cv := range connectionVariables {
				if c.ChildVariableID == cv.VariableID {
					connectionVariable = cv
				}
			}

			viewConnection := ConnectionVariable{
				VariableID:             connectionVariable.VariableID,
				Value:                  connectionVariable.Name,
				Path:                   c.Path,
				StructureType:          connectionVariable.StructureType,
				CreatifSpecialVariable: true,
			}

			b, err := json.Marshal(viewConnection)
			if err != nil {
				return nil, err
			}

			updatedValue, err := sjson.SetRawBytes(value, c.Path, b)
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

		for _, c := range conns {
			var connectionVariable QueryValueView
			for _, cv := range connectionVariables {
				if c.ChildVariableID == cv.VariableID {
					connectionVariable = cv
				}
			}

			updatedValue, err := sjson.SetRawBytes(value, c.Path, connectionVariable.Value)
			if err != nil {
				return nil, err
			}

			value = updatedValue
		}
	}

	return value, nil
}
