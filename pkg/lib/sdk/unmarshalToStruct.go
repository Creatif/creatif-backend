package sdk

import "encoding/json"

func UnmarshalToStruct[T any](source []byte) (T, error) {
	var target T

	if err := json.Unmarshal(source, &target); err != nil {
		return target, err
	}

	return target, nil
}
