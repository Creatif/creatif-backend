package sdk

import (
	"encoding/json"
)

func ErrorToResponseError(err error) map[string]string {
	var e map[string]string
	b, err := json.Marshal(err)
	if err != nil {
		return map[string]string{
			"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
		}
	}

	if err := json.Unmarshal(b, &e); err != nil {
		return map[string]string{
			"unrecoverable": "An internal validation error occurred. This should not happen. Please, submit a bug.",
		}
	}

	return e
}
