package sdk

import "encoding/json"

func CovertToGeneric(data []byte) ([]byte, error) {
	var model interface{}
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	b, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	return b, err
}
