package pagination

import (
	"encoding/base64"
	"encoding/json"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type cursor struct {
	nextID    string `json:"nextID"`
	prevID    string `json:"prevID"`
	field     string `json:"field"`
	orderBy   string `json:"orderBy"`
	direction string `json:"direction"`
	limit     int    `json:"limit"`
}

type Cursor interface {
	NextID() string
	PrevID() string
	Field() string
	OrderBy() string
	Direction() string
	Limit() int
}

func (c cursor) NextID() string {
	return c.nextID
}

func (c cursor) PrevID() string {
	return c.prevID
}

func (c cursor) Field() string {
	return c.field
}

func (c cursor) OrderBy() string {
	return c.orderBy
}

func (c cursor) Direction() string {
	return c.direction
}

func (c cursor) Limit() int {
	return c.limit
}

func CursorFromString(c string) (Cursor, error) {
	return decodeCursor(c)
}

func CursorFromData(nextId, prevId, field, orderBy, direction string, limit int) Cursor {
	return cursor{
		nextID:    nextId,
		prevID:    prevId,
		field:     field,
		orderBy:   orderBy,
		direction: direction,
		limit:     limit,
	}
}

func encodeCursor(c Cursor) (string, error) {
	serializedCursor, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	utf16Cursor, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), serializedCursor)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(utf16Cursor)
	return encodedCursor, nil
}

func decodeCursor(c string) (Cursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(c)
	if err != nil {
		return nil, err
	}

	utf8Cursor, _, err := transform.Bytes(unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(), decodedCursor)

	var cur cursor
	if err := json.Unmarshal(utf8Cursor, &cur); err != nil {
		return nil, err
	}

	return cur, nil
}
