package sdk

import (
	"bytes"
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func NewULID() (string, error) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// The slice should now contain random bytes instead of only zeroes.
	ms := ulid.Timestamp(time.Now())
	uid, err := ulid.New(ms, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}
