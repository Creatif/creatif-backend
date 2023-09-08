package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type cursor struct {
	ID      string `json:"id"`
	Field   string `json:"field"`
	OrderBy string `json:"orderBy"`
}

func newCursor(id, field, orderBy string) cursor {
	return cursor{
		ID:      id,
		Field:   field,
		OrderBy: orderBy,
	}
}

func resolveCursor(id, field, orderBy string) (string, error) {
	return encodeCursor(newCursor(id, field, orderBy))
}

func encodeCursor(c cursor) (string, error) {
	serializedCursor, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor, nil
}

func decodeCursor(c string) (cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return cursor{}, err
	}

	var cur cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return cursor{}, err
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
