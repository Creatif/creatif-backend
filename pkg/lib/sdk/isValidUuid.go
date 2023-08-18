package sdk

import "github.com/google/uuid"

func IsValidUuid(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
