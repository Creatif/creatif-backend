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

func (c *cursor) String() (string, error) {
	return encodeCursor(*c)
}

func CursorFromString(cur string) (*cursor, error) {
	return decodeCursor(cur)
}

func CursorFromData(id, field, orderBy string) *cursor {
	return &cursor{
		ID:      id,
		Field:   field,
		OrderBy: orderBy,
	}
}

func (c *cursor) encode() (string, error) {
	serializedCursor, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor, nil
}

func NewCursor(current string, id, field, orderBy string) (*cursor, error) {
	if current != "" {
		return CursorFromString(current)
	}

	return CursorFromData(id, field, orderBy), nil
}

func encodeCursor(c cursor) (string, error) {
	serializedCursor, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(serializedCursor)
	return encodedCursor, nil
}

func decodeCursor(c string) (*cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return nil, err
	}

	var cur cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return &cur, nil
}
