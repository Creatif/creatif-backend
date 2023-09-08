package sdk

import (
	"github.com/oklog/ulid/v2"
)

func IsValidUuid(u string) bool {
	_, err := ulid.Parse(u)
	return err == nil
}
