package create

import "encoding/json"

func decideToCreateNewActivity(lastWritten []byte, incoming []byte) (bool, error) {
	var lastWrittenDataQuery map[string]string
	var incomingDataQuery map[string]string
	if err := json.Unmarshal(lastWritten, &lastWrittenDataQuery); err != nil {
		return false, err
	}

	if err := json.Unmarshal(incoming, &incomingDataQuery); err != nil {
		return false, err
	}

	if lastWrittenDataQuery["type"] == "visit" && incomingDataQuery["type"] == "visit" {
		return false, nil
	}

	return false, nil
}
