package sdk

import "encoding/json"

func ConvertByUnmarshaling[T any](v interface{}, model T) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, model); err != nil {
		return err
	}

	return nil
}
