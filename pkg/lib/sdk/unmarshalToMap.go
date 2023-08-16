package sdk

import (
	"encoding/json"
)

func UnmarshalToMap[T string | []byte](v T) map[string]interface{} {
	b := []byte(v)
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		// log error but don't crash the program, log to slack also
	}

	return m
}
