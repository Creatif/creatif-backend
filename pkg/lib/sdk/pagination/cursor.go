package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type Cursor struct {
	ID      string `json:"id"`
	Field   string `json:"field"`
	OrderBy string `json:"orderBy"`
}

func NewCursor(id, field, orderBy string) Cursor {
	return Cursor{
		ID:      id,
		Field:   field,
		OrderBy: orderBy,
	}
}

func encodeCursor(c Cursor) (string, error) {
	serializedCursor, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor, nil
}

func decodeCursor(c string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return Cursor{}, err
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return Cursor{}, err
	}

	return cur, nil
}

func getPaginationOperator(direction string, sortOrder string) (string, string) {
	if direction == DIRECTION_FORWARD && sortOrder == "asc" {
		return ">", ""
	}
	if direction == DIRECTION_FORWARD && sortOrder == "desc" {
		return "<", ""
	}
	if direction == DIRECTION_BACKWARDS && sortOrder == "asc" {
		return "<", "desc"
	}
	if direction != DIRECTION_BACKWARDS && sortOrder == "desc" {
		return ">", "asc"
	}

	return "", ""
}
